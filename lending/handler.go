package lending

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/config"
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/internalerror"
	"github.com/radianhanggata/siesta-coding-test/lending-account-svc/model"
)

type Handler struct {
	Repository Repository
	Config     config.Repository
}

func SetupHandler(repository Repository, config config.Repository) Handler {
	return Handler{
		Repository: repository,
		Config:     config,
	}
}

func (h *Handler) InsertLending(c *fiber.Ctx) error {
	request := &InsertLendingRequest{}
	err := c.BodyParser(request)
	if err != nil {
		return c.Status(internalerror.ErrBadRequest.Code).JSON(internalerror.ErrBadRequest)
	}

	lendingConfig, err := h.Config.GetByID("lending")
	if err != nil && err == internalerror.ErrNotFound {
		return c.Status(ConfigNotFound.Code).JSON(ConfigNotFound)
	}

	fee := lendingConfig.Fee * 0.01 * request.Amount
	interest := lendingConfig.Interest * 0.01 * request.Amount * float64(request.Tenor)

	entity := model.Lending{
		Date:      time.Now(),
		Amount:    request.Amount,
		Tenor:     request.Tenor,
		Fee:       fee,
		Interest:  interest,
		AccountID: request.AccountID,
	}

	err = h.Repository.InsertLending(&entity)
	if err != nil {
		response := err.(*internalerror.Response)
		return c.Status(response.Code).JSON(response)
	}

	listRepayment := make([]*model.Repayment, 0)
	paymentDate := entity.Date.AddDate(0, 1, 0)
	for i := 0; i < request.Tenor; i++ {
		repayment := &model.Repayment{
			Date:      paymentDate,
			Fee:       fee,
			Interest:  interest / float64(entity.Tenor),
			Principal: entity.Amount / float64(entity.Tenor),
			AccountID: entity.AccountID,
			LendingID: entity.ID,
		}

		listRepayment = append(listRepayment, repayment)

		paymentDate = paymentDate.AddDate(0, 1, 0)
		fee = 0
	}

	err = h.Repository.InsertRepayment(listRepayment)
	if err != nil {
		response := err.(*internalerror.Response)
		return c.Status(response.Code).JSON(response)
	}

	return c.Status(201).Send(nil)
}

func calculateRepayment(accountID uint, listUnpaidRepayment []model.Repayment, config model.Config) []SimulateRepaymentResponse {
	listSimulateRepayment := make([]SimulateRepaymentResponse, 0)

	outstanding := float64(0)
	for _, unpaidRepayment := range listUnpaidRepayment {
		outstanding += unpaidRepayment.Principal
	}

	var month time.Time
	index := -1
	for i := 0; i < len(listUnpaidRepayment); i++ {
		if month.Month() != listUnpaidRepayment[i].Date.Month() {
			month = listUnpaidRepayment[i].Date
			listSimulateRepayment = append(listSimulateRepayment, SimulateRepaymentResponse{})
			if index > -1 {
				listSimulateRepayment[index].FeeStampDuty = getFeeStampDuty(outstanding, config)
				listSimulateRepayment[index].Tagihan += listSimulateRepayment[index].FeeStampDuty
				outstanding -= listSimulateRepayment[index].PokokYangDibayar
				listSimulateRepayment[index].Outstanding = math.Round(outstanding)
			}
			index++
		}

		listSimulateRepayment[index].Month = int(listUnpaidRepayment[i].Date.Month())
		listSimulateRepayment[index].Year = int(listUnpaidRepayment[i].Date.Year())
		listSimulateRepayment[index].Fee += listUnpaidRepayment[i].Fee
		listSimulateRepayment[index].Interest += listUnpaidRepayment[i].Interest
		listSimulateRepayment[index].PokokYangDibayar += math.Round(listUnpaidRepayment[i].Principal)
		listSimulateRepayment[index].Tagihan += math.Round(listUnpaidRepayment[i].Principal + listUnpaidRepayment[i].Interest + listUnpaidRepayment[i].Fee)
	}

	return listSimulateRepayment
}

func getFeeStampDuty(outstanding float64, config model.Config) float64 {
	if outstanding > config.OutstandingThreshold {
		return config.OutstandingFee
	}
	return 0
}

func (h *Handler) Simulate(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(internalerror.ErrBadRequest.Code).JSON(internalerror.ErrBadRequest)
	}

	q := c.Queries()
	amount, err := strconv.ParseFloat(q["amount"], 64)
	if err != nil {
		return c.Status(internalerror.ErrBadRequest.Code).JSON(internalerror.ErrBadRequest)
	}

	tenor, err := strconv.Atoi(q["tenor"])
	if err != nil {
		return c.Status(internalerror.ErrBadRequest.Code).JSON(internalerror.ErrBadRequest)
	}

	existingRepayment, err := h.Repository.GetUnpaidRepayment(uint(id))
	if err != nil {
		response := err.(*internalerror.Response)
		return c.Status(response.Code).JSON(response)
	}

	lendingConfig, err := h.Config.GetByID("lending")
	if err != nil && err == internalerror.ErrNotFound {
		return c.Status(ConfigNotFound.Code).JSON(ConfigNotFound)
	}

	fee := lendingConfig.Fee * 0.01 * amount
	interest := float64(0)
	if tenor > 1 {
		interest = lendingConfig.Interest * 0.01 * amount
	}
	principal := math.Round(amount / float64(tenor))

	addedRepayment := make([]model.Repayment, 0)
	paymentDate := time.Now().AddDate(0, 1, 0)
	for i := 0; i < tenor; i++ {
		repayment := model.Repayment{
			Date:      paymentDate,
			Fee:       fee,
			Interest:  interest,
			Principal: principal,
			AccountID: uint(id),
		}

		addedRepayment = append(addedRepayment, repayment)

		paymentDate = paymentDate.AddDate(0, 1, 0)
		fee = 0
	}

	addedRepayment = append(addedRepayment, existingRepayment...)

	sort.SliceStable(addedRepayment, func(i, j int) bool { return addedRepayment[i].Date.Before(addedRepayment[j].Date) })

	newRepayment := calculateRepayment(uint(id), addedRepayment, lendingConfig)

	outstanding := amount
	for i := 0; i < len(newRepayment); i++ {
		outstanding += newRepayment[i].PokokYangDibayar
	}

	feeStampDuty := float64(0)
	fee = 0
	interest = 0
	total := float64(0)
	for _, value := range newRepayment {
		feeStampDuty += value.FeeStampDuty
		fee += value.Fee
		interest += value.Interest
		total += value.Tagihan
	}

	response := &SimulateLendingResponse{
		Fee:          fee,
		FeeStampDuty: feeStampDuty,
		Interest:     interest,
		TotalPayment: total,
	}

	return c.Status(200).JSON(response)

}

func (h *Handler) GetRepayment(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(internalerror.ErrBadRequest.Code).JSON(internalerror.ErrBadRequest)
	}

	listUnpaidRepayment, err := h.Repository.GetUnpaidRepayment(uint(id))
	if err != nil {
		return c.Status(500).Send([]byte("error"))
	}

	lendingConfig, err := h.Config.GetByID("lending")
	if err != nil && err == internalerror.ErrNotFound {
		return c.Status(ConfigNotFound.Code).JSON(ConfigNotFound)
	}

	listSimulateRepayment := calculateRepayment(uint(id), listUnpaidRepayment, lendingConfig)
	if len(listSimulateRepayment) == 0 {
		return c.Status(RepaymentNotFound.Code).JSON(RepaymentNotFound)
	}

	return c.Status(200).JSON(listSimulateRepayment)
}
