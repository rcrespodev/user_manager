package repository

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"

type MessagesSchema struct {
	Messages []MessageSchema
}

type MessageSchema struct {
	Id              message.MessageId
	Pkg             string
	Text            string
	ClientErrorType message.ClientErrorType
}

func GetMessagesSourceData() *MessagesSchema {
	return &MessagesSchema{
		Messages: []MessageSchema{
			{
				Id:              1,
				Pkg:             "user",
				Text:            "user %v created successful",
				ClientErrorType: 0,
			},
			{
				Id:              4,
				Pkg:             "user",
				Text:            "value %v is invalid as %v attribute",
				ClientErrorType: 1,
			},
			{
				Id:              5,
				Pkg:             "user",
				Text:            "attribute %v are mandatory",
				ClientErrorType: 1,
			},
			{
				Id:              6,
				Pkg:             "user",
				Text:            "attribute %v can´t be greater than %v characters",
				ClientErrorType: 1,
			},
			{
				Id:              7,
				Pkg:             "user",
				Text:            "attribute %v can´t contain special characters (%v)",
				ClientErrorType: 1,
			},
			{
				Id:              8,
				Pkg:             "user",
				Text:            "password must be contain at least one special character like %$#&",
				ClientErrorType: 1,
			},
			{
				Id:              9,
				Pkg:             "user",
				Text:            "password must be contain at least one number",
				ClientErrorType: 1,
			},
			{
				Id:              10,
				Pkg:             "user",
				Text:            "attribute %v can´t be smaller than %v characters",
				ClientErrorType: 1,
			},
			{
				Id:              11,
				Pkg:             "user",
				Text:            "password must be contain at least one upper case",
				ClientErrorType: 1,
			},
			{
				Id:              12,
				Pkg:             "user",
				Text:            "password must be contain at least one lower case",
				ClientErrorType: 1,
			},
			{
				Id:              13,
				Pkg:             "user",
				Text:            "%v attribute dont´t must contain %v",
				ClientErrorType: 1,
			},
			{
				Id:              14,
				Pkg:             "user",
				Text:            "user with component: %v and value: %v already exists",
				ClientErrorType: 1,
			},
		},
	}
}
