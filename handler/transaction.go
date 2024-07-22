package handler

import (
	"maglo/database"
	"maglo/model"

	"github.com/gofiber/fiber/v2"
)

// GetAllTransactions query all transactions
func GetAllTransactions(c *fiber.Ctx) error {
	db := database.DB
	var transactions []model.Transaction
	db.Find(&transactions)
	return c.JSON(fiber.Map{"status": "success", "message": "All transactions", "data": transactions})
}

// GetTransaction query transaction
func GetTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var transaction model.Transaction
	db.Find(&transaction, id)
	if transaction.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No transaction found with ID", "data": nil})

	}
	return c.JSON(fiber.Map{"status": "success", "message": "Transaction found", "data": transaction})
}

// CreateTransaction new transaction
func CreateTransaction(c *fiber.Ctx) error {
	db := database.DB
	transaction := new(model.Transaction)
	if err := c.BodyParser(transaction); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create transaction", "data": err})
	}
	db.Create(&transaction)
	return c.JSON(fiber.Map{"status": "success", "message": "Created transaction", "data": transaction})
}

// DeleteTransaction delete transaction
func DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var transaction model.Transaction
	db.First(&transaction, id)
	if transaction.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No transaction found with ID", "data": nil})

	}
	db.Delete(&transaction)
	return c.JSON(fiber.Map{"status": "success", "message": "Transaction successfully deleted", "data": nil})
}