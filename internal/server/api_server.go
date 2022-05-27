package server

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/ArtyomArtamonov/msg/internal/model"
	"github.com/ArtyomArtamonov/msg/internal/repository"
	pb "github.com/ArtyomArtamonov/msg/internal/server/msg-proto"
	"github.com/ArtyomArtamonov/msg/internal/service"
	"github.com/ArtyomArtamonov/msg/internal/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ApiServer struct {
	pb.UnimplementedApiServiceServer

	jwtManager   service.JWTManagerProtol
	roomStore    repository.RoomStore
	messageStore repository.MessageStore
	amqpManager  service.AMQPProducer
}

func NewApiServer(jwtManager service.JWTManagerProtol, roomStore repository.RoomStore, messageStore repository.MessageStore, amqpManager service.AMQPProducer) *ApiServer {
	return &ApiServer{
		jwtManager:   jwtManager,
		roomStore:    roomStore,
		messageStore: messageStore,
		amqpManager:  amqpManager,
	}
}

func (s *ApiServer) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomStatus, error) {
	claims, err := s.jwtManager.GetAndVerifyClaims(ctx)
	if err != nil {
		return nil, err
	}

	usersInRoom := map[string]bool{}
	usersInRoom[claims.Id] = true
	for _, userId := range req.UserIds {
		usersInRoom[userId] = true
	}

	usersInRoomUUIDs := []uuid.UUID{}
	for userId := range usersInRoom {
		id, err := uuid.Parse(userId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "could not parse uuid: %v", err)
		}
		usersInRoomUUIDs = append(usersInRoomUUIDs, id)
	}

	newRoom := model.NewRoom(req.Name, false, usersInRoomUUIDs...)
	err = s.roomStore.Add(newRoom)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create room: %v", err)
	}

	userIdsResponse := []string{}
	for _, id := range usersInRoomUUIDs {
		userIdsResponse = append(userIdsResponse, id.String())
	}

	createRoomStatus := &pb.CreateRoomStatus{
		RoomId: newRoom.Id.String(),
		Name:   newRoom.Name,
		Users:  userIdsResponse,
	}

	return createRoomStatus, nil
}

func (s *ApiServer) ListRooms(ctx context.Context, req *pb.ListRoomsRequest) (*pb.ListRoomsResponse, error) {
	if req.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "page_size cannot be bigger than 100")
	}

	claims, err := s.jwtManager.GetAndVerifyClaims(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot parse uuid: %v", err)
	}

	var rooms []model.Room
	if req.NextToken == nil {
		rooms, err = s.roomStore.ListRoomsFirst(userId, int(req.PageSize))
		if err != nil {
			return nil, err
		}
	} else {
		lastMessageTime, e := decodePageToken(req.NextToken.Value)
		if e != nil {
			return nil, status.Errorf(codes.InvalidArgument, "cannot parse next token: %v", e)
		}
		rooms, err = s.roomStore.ListRooms(userId, *lastMessageTime, int(req.PageSize))
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot get rooms: %v", err)
	}

	var nextToken string
	if len(rooms) < int(req.PageSize) {
		nextToken = ""
	} else {
		lastRoom := rooms[len(rooms)-1]
		nextToken = encodePageToken(lastRoom.LastMessageTime)
	}

	pbRooms := []*pb.Room{}
	for _, room := range rooms {
		pbRooms = append(pbRooms, room.PbRoom())
	}

	var listRoomsResponse *pb.ListRoomsResponse
	if len(nextToken) == 0 {
		listRoomsResponse = &pb.ListRoomsResponse{
			NextToken: nil,
			Rooms:     pbRooms,
		}
	} else {
		listRoomsResponse = &pb.ListRoomsResponse{
			NextToken: &wrapperspb.StringValue{Value: nextToken},
			Rooms:     pbRooms,
		}
	}

	return listRoomsResponse, nil
}

func (s *ApiServer) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	if req.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "page_size cannot be bigger than 100")
	}

	claims, err := s.jwtManager.GetAndVerifyClaims(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot parse uuid: %v", err)
	}

	chatId, err := uuid.Parse(req.ChatId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "could not parse uuid")
	}

	roomUserIds, err := s.roomStore.UsersInRoom(chatId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}

	if !utils.ArrayContains(roomUserIds, userId) {
		return nil, status.Error(codes.PermissionDenied, "")
	}

	var messages []model.Message
	if req.NextToken == nil {
		messages, err = s.messageStore.ListMessagesFirst(chatId, int(req.PageSize))
		if err != nil {
			return nil, err
		}
	} else {
		lastMessageTime, e := decodePageToken(req.NextToken.Value)
		if e != nil {
			return nil, status.Errorf(codes.InvalidArgument, "cannot parse next token: %v", e)
		}
		messages, err = s.messageStore.ListMessages(chatId, *lastMessageTime, int(req.PageSize))
	}

	var nextToken string
	if len(messages) < int(req.PageSize) {
		nextToken = ""
	} else {
		lastMessage := messages[len(messages)-1]
		nextToken = encodePageToken(lastMessage.CreatedAt)
	}

	pbMessages := []*pb.Message{}
	for _, message := range messages {
		pbMessages = append(pbMessages, message.ToPbMessage())
	}

	var listMessagesResponse *pb.ListMessagesResponse
	if len(nextToken) == 0 {
		listMessagesResponse = &pb.ListMessagesResponse{
			NextToken: nil,
			Messages:  pbMessages,
		}
	} else {
		listMessagesResponse = &pb.ListMessagesResponse{
			NextToken: &wrapperspb.StringValue{Value: nextToken},
			Messages:  pbMessages,
		}
	}

	return listMessagesResponse, nil
}

func (s *ApiServer) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	claims, err := s.jwtManager.GetAndVerifyClaims(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot parse uuid: %v", err)
	}

	senderId, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "could not parse uuid")
	}

	room, ok := req.Recipient.(*pb.MessageRequest_RoomId)

	if !ok {
		// Dialog room is not created, first message is about to send

		userId := req.Recipient.(*pb.MessageRequest_UserId)

		recipientId, err := uuid.Parse(userId.UserId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "could not parse uuid")
		}

		room := model.NewRoom("", true, senderId, recipientId)
		message := model.NewMessage(senderId, uuid.Nil, req.Message)
		roomResponse, err := s.roomStore.AddAndSendMessage(room, message)
		if err != nil && status.Code(err) == codes.AlreadyExists {
			message.RoomId = roomResponse.Id

			err := s.messageStore.SendMessage(message)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "could not create room or send message: %v", err)
			}
		}

		messageDelivery := &pb.MessageDelivery{
			Message: message.ToPbMessage(),
			UserIds: []string{recipientId.String()},
		}
		err = s.amqpManager.Produce(messageDelivery)
		if err != nil {
			logrus.Errorf("could not send message by amqp: %v", err)
		}

		response := &pb.MessageResponse{
			RoomId:  roomResponse.PbRoom().Id,
			Message: message.ToPbMessage(),
		}
		return response, nil
	} else {
		// Sending message to already existing room

		roomId, err := uuid.Parse(room.RoomId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "could not parse uuid")
		}

		roomUserIds, err := s.roomStore.UsersInRoom(roomId)
		if err != nil {
			return nil, status.Error(codes.NotFound, "")
		}

		if !utils.ArrayContains(roomUserIds, userId) {
			return nil, status.Error(codes.PermissionDenied, "")
		}

		message := model.NewMessage(senderId, roomId, req.Message)
		err = s.messageStore.SendMessage(message)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not send message: %v", err)
		}
		var recipientUserIds []string
		for _, roomUserId := range roomUserIds {
			recipientUserIds = append(recipientUserIds, roomUserId.String())
		}

		messageDelivery := &pb.MessageDelivery{
			Message: message.ToPbMessage(),
			UserIds: recipientUserIds,
		}
		err = s.amqpManager.Produce(messageDelivery)
		if err != nil {
			logrus.Errorf("could not send message by amqp: %v", err)
		}

		response := &pb.MessageResponse{
			RoomId:  room.RoomId,
			Message: message.ToPbMessage(),
		}
		return response, nil
	}
}

func decodePageToken(token string) (*time.Time, error) {
	buf, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, string(buf))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func encodePageToken(lastMessageDate time.Time) string {
	return base64.StdEncoding.EncodeToString([]byte(lastMessageDate.Format(time.RFC3339)))
}
