package depositvalidations

import (
	"context"
	"errors"
	"fmt"
	"heypay-cash-in-server/internal/rest"
	"heypay-cash-in-server/internal/services/accountsservice"
	"heypay-cash-in-server/internal/services/databaseservice"
	"heypay-cash-in-server/internal/services/databaseservice/transferpaymentrequestsservice"
	"heypay-cash-in-server/internal/services/databaseservice/usersservice"
	"heypay-cash-in-server/internal/services/firebaseservice"
	"heypay-cash-in-server/internal/services/notificationservice"
	"heypay-cash-in-server/models/account"
	"heypay-cash-in-server/models/transferpaymentrequest"
	"heypay-cash-in-server/models/user"
	"heypay-cash-in-server/tests"
	"heypay-cash-in-server/tests/utils/httpclient"
	"heypay-cash-in-server/utils"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

const Endpoint = "deposit-validations"
const Method = "POST"

func init() {
	utils.DisableLogging()
	go setupServer()
	time.Sleep(5000)
}

func setupServer() {
	utils.Info(nil, "Starting the service...")
	router := rest.Router()
	utils.Info(nil, "The service is ready to listen and serve.")
	log.Fatal(http.ListenAndServe(":8000", router))
}

var _ = Describe(fmt.Sprintf("/%s", Endpoint), func() {
	BeforeEach(func() {
		// Reset each mocked service
		transferpaymentrequestsservice.FindPendingTransferPaymentRequest = func(context.Context, int64, string, string, string, time.Time, *firestore.Transaction) (*transferpaymentrequest.TransferPaymentRequest, error) {
			return nil, nil
		}
		transferpaymentrequestsservice.UpdateDocument = func(context.Context, string, []firestore.Update, databaseservice.UpdateDocumentConfig) error {
			return nil
		}
		usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
			return nil, nil
		}
		accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
			return nil, nil
		}
		notificationservice.SendPushNotification = func(context.Context, string, notificationservice.Payload) error {
			return nil
		}
		firebaseservice.CheckMaxDailyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
			return firebaseservice.RuleCheckResult{}, nil
		}
		firebaseservice.CheckMaxMonthlyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
			return firebaseservice.RuleCheckResult{}, nil
		}
		firebaseservice.CheckMaxDailyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
			return firebaseservice.RuleCheckResult{}, nil
		}
		firebaseservice.CheckMaxMonthlyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
			return firebaseservice.RuleCheckResult{}, nil
		}
	})
	Describe(fmt.Sprintf("%s", Method), func() {
		Context("When sending invalid data (missing one or more mandatory values).", func() {
			for i := 0; i < len(InvalidData)-1; i++ {
				Describe(fmt.Sprintf("Sending %v in body.", InvalidData[i]), func() {
					It("should respond with status 400 and code \"invalid-parameters\" in the body.", func() {
						var response Response
						err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), InvalidData[i], httpclient.HttpConfig{}, &response)
						if err != nil {
							Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
							return
						}
						Expect(response).NotTo(BeNil())
						Expect(response.Status).To(Equal(http.StatusBadRequest))
						Expect(response.Code).To(Equal("invalid-parameters"))
					})
				})
			}
		})
		Context("When sending valid data, but origin rut its on banned origin rut list.", func() {
			var bannedOriginRutBody = map[string]interface{}{}
			for k, v := range ValidData[0].Data {
				bannedOriginRutBody[k] = v
			}
			bannedOriginRutBody["originRut"] = "0000176988703"
			Describe(fmt.Sprintf("Sending %v in body.", bannedOriginRutBody), func() {
				It("should respond with status 200 and code \"origin-rut-banned\" in the body.", func() {
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), bannedOriginRutBody, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("origin-rut-banned"))
				})
			})
		})
		Context("When sending valid data, but triggering firestore errors.", func() {
			Describe(fmt.Sprintf("Sending %v in body.", ValidData[0].Data), func() {
				It("should respond with status 500 and code \"05\" in the body when finding pending transfer payment requests fails.", func() {
					transferpaymentrequestsservice.FindPendingTransferPaymentRequest = func(context.Context, int64, string, string, string, time.Time, *firestore.Transaction) (*transferpaymentrequest.TransferPaymentRequest, error) {
						return nil, errors.New("random firestore error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("05"))
				})
				It("should respond with status 500 and code \"06\" in the body when updating transfer payment requests fails.", func() {
					transferpaymentrequestsservice.FindPendingTransferPaymentRequest = func(context.Context, int64, string, string, string, time.Time, *firestore.Transaction) (*transferpaymentrequest.TransferPaymentRequest, error) {
						return TransferPaymentRequestDocumentMock, nil
					}
					transferpaymentrequestsservice.UpdateDocument = func(context.Context, string, []firestore.Update, databaseservice.UpdateDocumentConfig) error {
						return errors.New("random firestore error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("06"))
				})
				It("should respond with status 500 and code \"06\" in the body when transfer payment requests transaction fails.", func() {
					// TODO: find a way to trigger a transaction error, probably we will need to encapsulate firestore implementation and override its behaviour like other database service methods
					// InternalServerError - 06 - Error in transfer payment requests transaction - reason: %v
				})
				It("should respond with status 500 and code \"07\" in the body when finding user fails.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return nil, errors.New("random firestore error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("07"))
				})
				It("should respond with status 500 and code \"10\" in the body when checking max daily bank transfer deposit amount fails.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return UserAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, errors.New("random firestore error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("10"))
				})
				It("should respond with status 500 and code \"12\" in the body when checking max monthly bank transfer deposit amount fails.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return UserAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{Check: true}, nil
					}
					firebaseservice.CheckMaxMonthlyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, errors.New("random firestore error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("12"))
				})
				It("should respond with status 500 and code \"14\" in the body when checking max daily bank transfer payment amount fails.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return CommerceAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, errors.New("random firestore error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("14"))
				})
				It("should respond with status 500 and code \"15\" in the body when checking max monthly bank transfer payment amount fails.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return CommerceAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{Check: true}, nil
					}
					firebaseservice.CheckMaxMonthlyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, errors.New("random firestore error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("15"))
				})
			})
		})
		Context("When sending valid data, but user does not exist.", func() {
			Describe(fmt.Sprintf("Sending %v in body.", ValidData[0].Data), func() {
				It("should respond with status 200 and code \"user-not-found\".", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{}, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("user-not-found"))
				})
			})
		})
		Context("When sending valid data, but triggering accounts engine service errors.", func() {
			Describe(fmt.Sprintf("Sending %v in body.", ValidData[0].Data), func() {
				It("should respond with status 500 and code \"08\".", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return nil, errors.New("random accounts engine error")
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(http.StatusInternalServerError))
					Expect(response.Code).To(Equal("08"))
				})
			})
		})
		Context("When sending valid data, but the account status is not active.", func() {
			Describe(fmt.Sprintf("Sending %v in body.", ValidData[0].Data), func() {
				It("should respond with status 200 and code \"account-blocked\".", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return BlockedAccountDocumentMock, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("account-blocked"))
				})
			})
		})
		Context("When sending valid data, but the account would be over its max balance.", func() {
			Describe(fmt.Sprintf("Sending %v in body.", ValidData[0].Data), func() {
				It("should respond with status 200 and code \"over-max-balance\".", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return OverMaxBalanceAccountDocumentMock, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("over-max-balance"))
				})
			})
		})
		Context("When sending valid data, but the account does not meet transactional limit rules.", func() {
			Describe(fmt.Sprintf("Sending %v in body.", ValidData[0].Data), func() {
				It("should respond with status 200 and code \"over-daily-deposit-amount\", when the account would be over its daily deposit amount.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return UserAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("over-daily-deposit-amount"))
				})
				It("should respond with status 200 and code \"over-monthly-deposit-amount\", when the account would be over its monthly deposit amount.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return UserAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{Check: true}, nil
					}
					firebaseservice.CheckMaxMonthlyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("over-monthly-deposit-amount"))
				})
				It("should respond with status 200 and code \"over-daily-payment-amount\", when the account would be over its daily deposit payment amount.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return CommerceAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("over-daily-deposit-payment-amount"))
				})
				It("should respond with status 200 and code \"over-monthly-payment-amount\", when the account would be over its monthly deposit payment amount.", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return CommerceAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{Check: true}, nil
					}
					firebaseservice.CheckMaxMonthlyBankTransferPaymentAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{}, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("over-monthly-deposit-payment-amount"))
				})
			})
		})
		Context("When sending valid data.", func() {
			Describe(fmt.Sprintf("Sending %v in body.", ValidData[0].Data), func() {
				It("should respond with status 200 and code \"valid\".", func() {
					usersservice.GetDocumentsList = func(context.Context, map[string]databaseservice.QueryFieldConfig, *databaseservice.QueryOptions) ([]user.User, error) {
						return []user.User{UserDocumentMock}, nil
					}
					accountsservice.GetAccount = func(context.Context, string) (*account.Account, error) {
						return UserAccountDocumentMock, nil
					}
					firebaseservice.CheckMaxDailyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{Check: true}, nil
					}
					firebaseservice.CheckMaxMonthlyBankTransferDepositAmount = func(context.Context, int64, string, string) (firebaseservice.RuleCheckResult, error) {
						return firebaseservice.RuleCheckResult{Check: true}, nil
					}
					var response Response
					err := httpclient.Post(nil, fmt.Sprintf("%s/%s", tests.ApiUrl, Endpoint), ValidData[0].Data, httpclient.HttpConfig{}, &response)
					if err != nil {
						Expect(err).To(BeNil(), fmt.Sprintf("Test failed while calling endpoint - reason: %v", err))
						return
					}
					Expect(response).NotTo(BeNil())
					Expect(response.Status).To(Equal(0)) // its the equivalent to 200
					Expect(response.Code).To(Equal("valid"))
				})
			})
		})
	})
})
