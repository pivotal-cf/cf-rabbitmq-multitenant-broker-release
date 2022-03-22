package rabbithutch_test

import (
	"errors"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "rabbitmq-service-broker/rabbithutch"
	"rabbitmq-service-broker/rabbithutch/fakes"

	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

var _ = Describe("Deleting the user", func() {
	var (
		rabbitClient *fakes.FakeAPIClient
		rabbithutch  RabbitHutch
		body         *fakeBody
	)

	BeforeEach(func() {
		rabbitClient = new(fakes.FakeAPIClient)
		rabbithutch = New(rabbitClient)
		body = &fakeBody{}
	})

	Describe("DeleteUserAndConnections()", func() {
		Context("delete returns a response", func() {
			AfterEach(func() {
				Expect(body.Closed).To(BeTrue())
			})
			It("deletes the user", func() {
				rabbitClient.DeleteUserReturns(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)

				err := rabbithutch.DeleteUserAndConnections("fake-user")
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbitClient.DeleteUserCallCount()).To(Equal(1))
				Expect(rabbitClient.DeleteUserArgsForCall(0)).To(Equal("fake-user"))
			})

			It("closes all connections for the specified user", func() {
				connections := []rabbithole.ConnectionInfo{
					rabbithole.ConnectionInfo{
						Name: "Connection 1",
						User: "fake-user",
					},
					rabbithole.ConnectionInfo{
						Name: "Connection 2",
						User: "fake-user",
					},
					rabbithole.ConnectionInfo{
						Name: "Connection 3",
						User: "not-fake-user",
					},
				}

				rabbitClient.ListConnectionsReturns(connections, nil)
				rabbitClient.DeleteUserReturns(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)
				connectionbody1 := &fakeBody{}
				connectionbody2 := &fakeBody{}
				rabbitClient.CloseConnectionReturnsOnCall(0, &http.Response{StatusCode: http.StatusOK, Body: connectionbody1}, nil)
				rabbitClient.CloseConnectionReturnsOnCall(1, &http.Response{StatusCode: http.StatusOK, Body: connectionbody2}, nil)

				err := rabbithutch.DeleteUserAndConnections("fake-user")
				Expect(err).NotTo(HaveOccurred())

				Expect(rabbitClient.ListConnectionsCallCount()).To(Equal(1))
				Expect(rabbitClient.CloseConnectionCallCount()).To(Equal(2))
				Expect(rabbitClient.CloseConnectionArgsForCall(0)).To(Equal("Connection 1"))
				Expect(rabbitClient.CloseConnectionArgsForCall(1)).To(Equal("Connection 2"))
				Expect(connectionbody1.Closed).To(BeTrue())
				Expect(connectionbody2.Closed).To(BeTrue())
			})

		})

		Context("DeleteUser returns an error", func() {
			It("still attempts to delete the connections", func() {
				connections := []rabbithole.ConnectionInfo{
					rabbithole.ConnectionInfo{
						Name: "Connection 1",
						User: "fake-user",
					},
				}
				rabbitClient.ListConnectionsReturns(connections, nil)
				rabbitClient.DeleteUserReturns(nil, errors.New("fake user error"))
				connectionbody1 := &fakeBody{}
				rabbitClient.CloseConnectionReturnsOnCall(0, &http.Response{StatusCode: http.StatusOK, Body: connectionbody1}, nil)

				respErr := rabbithutch.DeleteUserAndConnections("fake-user")

				Expect(respErr).To(MatchError("fake user error"))
				Expect(rabbitClient.CloseConnectionCallCount()).To(Equal(1))
				Expect(rabbitClient.DeleteUserCallCount()).To(Equal(1))
				Expect(rabbitClient.DeleteUserArgsForCall(0)).To(Equal("fake-user"))
				Expect(connectionbody1.Closed).To(BeTrue())
			})
		})
	})

	Describe("DeleteUser()", func() {
		It("deletes the user", func() {
			rabbitClient.DeleteUserReturns(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)
			err := rabbithutch.DeleteUser("fake-user")
			Expect(err).NotTo(HaveOccurred())

			Expect(rabbitClient.DeleteUserCallCount()).To(Equal(1))
			Expect(rabbitClient.DeleteUserArgsForCall(0)).To(Equal("fake-user"))
		})

		It("returns an error if it cannot delete the user", func() {
			rabbitClient.DeleteUserReturns(&http.Response{Body: body}, errors.New("fake user error"))

			respErr := rabbithutch.DeleteUser("fake-user")

			Expect(respErr).To(MatchError("fake user error"))
			Expect(rabbitClient.DeleteUserCallCount()).To(Equal(1))
			Expect(rabbitClient.DeleteUserArgsForCall(0)).To(Equal("fake-user"))
		})
	})
})
