package commands

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mmctl/printer"

	"github.com/spf13/cobra"
)

func (s *MmctlUnitTestSuite) TestSearchUserCmd() {
	s.Run("Search for an existing user", func() {
		emailArg := "example@example.com"
		mockUser := model.User{Username: "ExampleUser", Email: emailArg}

		s.client.
			EXPECT().
			GetUserByEmail(emailArg, "").
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		err := searchUserCmdF(s.client, &cobra.Command{}, []string{emailArg})
		s.Require().Nil(err)
		s.Require().Equal(&mockUser, printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Search for a nonexistent user", func() {
		printer.Clean()
		arg := "example@example.com"

		s.client.
			EXPECT().
			GetUserByEmail(arg, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetUserByUsername(arg, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			GetUser(arg, "").
			Return(nil, &model.Response{Error: nil}).
			Times(1)

		err := searchUserCmdF(s.client, &cobra.Command{}, []string{arg})
		s.Require().Nil(err)
		s.Require().Len(printer.GetLines(), 0)
		s.Require().Equal("Unable to find user 'example@example.com'", printer.GetErrorLines()[0])
	})
}

func (s *MmctlUnitTestSuite) TestUserCreateCmd() {
	mockUser := model.User{
		Username: "username",
		Password: "password",
		Email:    "email",
	}

	s.Run("Create user with email missing", func() {
		printer.Clean()

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Email is required: flag accessed but not defined: email", error.Error())
	})

	s.Run("Create user with username missing", func() {
		printer.Clean()

		command := cobra.Command{}
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Username is required: flag accessed but not defined: username", error.Error())
	})

	s.Run("Create user with password missing", func() {
		printer.Clean()

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Password is required: flag accessed but not defined: password", error.Error())
	})

	s.Run("Create a regular user", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Nil(error)
		s.Require().Equal(&mockUser, printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Create a regular user with client returning error", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: &model.AppError{Message: "Remote error"}}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Unable to create user. Error: : Remote error, ", error.Error())
	})

	s.Run("Create a sysAdmin user", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			UpdateUserRoles(mockUser.Id, "system_user system_admin").
			Return(true, &model.Response{Error: nil}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")
		command.Flags().Bool("system_admin", true, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Nil(error)
		s.Require().Equal(&mockUser, printer.GetLines()[0])
		s.Require().Len(printer.GetErrorLines(), 0)
	})

	s.Run("Create a sysAdmin user with client returning error", func() {
		printer.Clean()

		s.client.
			EXPECT().
			CreateUser(&mockUser).
			Return(&mockUser, &model.Response{Error: nil}).
			Times(1)

		s.client.
			EXPECT().
			UpdateUserRoles(mockUser.Id, "system_user system_admin").
			Return(false, &model.Response{Error: &model.AppError{Message: "Remote error"}}).
			Times(1)

		command := cobra.Command{}
		command.Flags().String("username", mockUser.Username, "")
		command.Flags().String("email", mockUser.Email, "")
		command.Flags().String("password", mockUser.Password, "")
		command.Flags().Bool("system_admin", true, "")

		error := userCreateCmdF(s.client, &command, []string{})

		s.Require().Equal("Unable to update user roles. Error: : Remote error, ", error.Error())
	})
}