package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cbqa/backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ─── HELPERS ─────────────────────────────────────────

func getID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err
}

func getUserID(c *gin.Context) uint {
	id, _ := c.Get("user_id")
	return id.(uint)
}

type PaginationQuery struct {
	Page  int    `form:"page,default=1"`
	Limit int    `form:"limit,default=10"`
	Q     string `form:"q"`
}

func paginate(q PaginationQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (q.Page - 1) * q.Limit
		return db.Offset(offset).Limit(q.Limit)
	}
}

// ─── DASHBOARD ───────────────────────────────────────

type DashboardHandler struct{ db *gorm.DB }

func NewDashboardHandler(db *gorm.DB) *DashboardHandler { return &DashboardHandler{db: db} }

func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := getUserID(c)
	var stats struct {
		OpenTasks           int64   `json:"open_tasks"`
		OpenProjects        int64   `json:"open_projects"`
		CompletedProjects   int64   `json:"completed_projects"`
		HoldProjects        int64   `json:"hold_projects"`
		TotalClients        int64   `json:"total_clients"`
		TotalLeads          int64   `json:"total_leads"`
		TotalMembers        int64   `json:"total_members"`
		DueAmount           float64 `json:"due_amount"`
		TotalIncome         float64 `json:"total_income"`
		TotalExpenses       float64 `json:"total_expenses"`
		TasksTodo           int64   `json:"tasks_todo"`
		TasksInProgress     int64   `json:"tasks_in_progress"`
		TasksDone           int64   `json:"tasks_done"`
		TasksExpired        int64   `json:"tasks_expired"`
		OverdueAmount       float64 `json:"overdue_amount"`
		NotPaidAmount       float64 `json:"not_paid_amount"`
		PartiallyPaidAmount float64 `json:"partially_paid_amount"`
		FullyPaidAmount     float64 `json:"fully_paid_amount"`
		DraftAmount         float64 `json:"draft_amount"`
		TotalInvoiced       float64 `json:"total_invoiced"`
		ClockedInCount      int64   `json:"clocked_in_count"`
		OnLeaveToday        int64   `json:"on_leave_today"`
	}
	today := time.Now().Format("2006-01-02")
	h.db.Model(&models.Task{}).Where("assigned_to_id = ? AND status != 'done'", userID).Count(&stats.OpenTasks)
	h.db.Model(&models.Project{}).Where("status = 'open'").Count(&stats.OpenProjects)
	h.db.Model(&models.Project{}).Where("status = 'completed'").Count(&stats.CompletedProjects)
	h.db.Model(&models.Project{}).Where("status = 'hold'").Count(&stats.HoldProjects)
	h.db.Model(&models.Client{}).Count(&stats.TotalClients)
	h.db.Model(&models.Lead{}).Count(&stats.TotalLeads)
	h.db.Model(&models.User{}).Where("is_active = true").Count(&stats.TotalMembers)
	h.db.Model(&models.Invoice{}).Select("COALESCE(SUM(due_amount),0)").Scan(&stats.DueAmount)
	h.db.Model(&models.Payment{}).Select("COALESCE(SUM(amount),0)").Scan(&stats.TotalIncome)
	h.db.Model(&models.Expense{}).Select("COALESCE(SUM(total),0)").Scan(&stats.TotalExpenses)
	h.db.Model(&models.Task{}).Where("status = 'todo'").Count(&stats.TasksTodo)
	h.db.Model(&models.Task{}).Where("status = 'in_progress'").Count(&stats.TasksInProgress)
	h.db.Model(&models.Task{}).Where("status = 'done'").Count(&stats.TasksDone)
	h.db.Model(&models.Task{}).Where("status = 'expired'").Count(&stats.TasksExpired)
	h.db.Model(&models.Invoice{}).Where("status = 'overdue'").Select("COALESCE(SUM(due_amount),0)").Scan(&stats.OverdueAmount)
	h.db.Model(&models.Invoice{}).Where("status = 'not_paid'").Select("COALESCE(SUM(total_amount),0)").Scan(&stats.NotPaidAmount)
	h.db.Model(&models.Invoice{}).Where("status = 'partially_paid'").Select("COALESCE(SUM(due_amount),0)").Scan(&stats.PartiallyPaidAmount)
	h.db.Model(&models.Invoice{}).Where("status = 'fully_paid'").Select("COALESCE(SUM(total_amount),0)").Scan(&stats.FullyPaidAmount)
	h.db.Model(&models.Invoice{}).Where("status = 'draft'").Select("COALESCE(SUM(total_amount),0)").Scan(&stats.DraftAmount)
	h.db.Model(&models.Invoice{}).Select("COALESCE(SUM(total_amount),0)").Scan(&stats.TotalInvoiced)
	h.db.Model(&models.User{}).Where("clocked_in = true AND is_active = true").Count(&stats.ClockedInCount)
	h.db.Model(&models.Leave{}).Where("status = 'approved' AND DATE(start_date) <= ? AND DATE(end_date) >= ?", today, today).Count(&stats.OnLeaveToday)
	c.JSON(http.StatusOK, stats)
}

// ─── CLIENT ──────────────────────────────────────────

type ClientHandler struct{ db *gorm.DB }

func NewClientHandler(db *gorm.DB) *ClientHandler { return &ClientHandler{db: db} }

func (h *ClientHandler) ListAllContacts(c *gin.Context) {
	var contacts []models.Contact
	h.db.Preload("Client").Find(&contacts)
	c.JSON(http.StatusOK, gin.H{"data": contacts, "total": len(contacts)})
}

func (h *ClientHandler) List(c *gin.Context) {
	var q PaginationQuery
	c.ShouldBindQuery(&q)
	var clients []models.Client
	var total int64
	query := h.db.Model(&models.Client{}).Preload("Owner").Preload("Labels")
	if q.Q != "" {
		query = query.Where("name ILIKE ?", "%"+q.Q+"%")
	}
	query.Count(&total)
	query.Scopes(paginate(q)).Find(&clients)
	c.JSON(http.StatusOK, gin.H{"data": clients, "total": total, "page": q.Page, "limit": q.Limit})
}

func (h *ClientHandler) Create(c *gin.Context) {
	var client models.Client
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	if client.OwnerID == 0 {
		client.OwnerID = getUserID(c)
	}
	if err := h.db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusCreated, client)
}

func (h *ClientHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var client models.Client
	if err := h.db.Preload("Owner").Preload("Contacts").Preload("Labels").First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"}); return
	}
	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var client models.Client
	if err := h.db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"}); return
	}
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Save(&client)
	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Client{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *ClientHandler) GetContacts(c *gin.Context) {
	id, _ := getID(c)
	var contacts []models.Contact
	h.db.Where("client_id = ?", id).Find(&contacts)
	c.JSON(http.StatusOK, gin.H{"data": contacts})
}

func (h *ClientHandler) AddContact(c *gin.Context) {
	id, _ := getID(c)
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	contact.ClientID = id
	h.db.Create(&contact)
	c.JSON(http.StatusCreated, contact)
}

func (h *ClientHandler) UpdateContact(c *gin.Context) {
	contactID, _ := strconv.ParseUint(c.Param("contactId"), 10, 64)
	var contact models.Contact
	if err := h.db.First(&contact, contactID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"}); return
	}
	c.ShouldBindJSON(&contact)
	h.db.Save(&contact)
	c.JSON(http.StatusOK, contact)
}

func (h *ClientHandler) DeleteContact(c *gin.Context) {
	contactID, _ := strconv.ParseUint(c.Param("contactId"), 10, 64)
	h.db.Delete(&models.Contact{}, contactID)
	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted"})
}

func (h *ClientHandler) GetProjects(c *gin.Context) {
	id, _ := getID(c)
	var projects []models.Project
	h.db.Where("client_id = ?", id).Find(&projects)
	c.JSON(http.StatusOK, gin.H{"data": projects})
}

func (h *ClientHandler) GetInvoices(c *gin.Context) {
	id, _ := getID(c)
	var invoices []models.Invoice
	h.db.Where("client_id = ?", id).Find(&invoices)
	c.JSON(http.StatusOK, gin.H{"data": invoices})
}

// ─── PROJECT ─────────────────────────────────────────

type ProjectHandler struct{ db *gorm.DB }

func NewProjectHandler(db *gorm.DB) *ProjectHandler { return &ProjectHandler{db: db} }

func (h *ProjectHandler) List(c *gin.Context) {
	var q PaginationQuery
	c.ShouldBindQuery(&q)
	var projects []models.Project
	var total int64
	query := h.db.Model(&models.Project{}).Preload("Client").Preload("Labels")
	if q.Q != "" {
		query = query.Where("title ILIKE ?", "%"+q.Q+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	query.Scopes(paginate(q)).Find(&projects)
	c.JSON(http.StatusOK, gin.H{"data": projects, "total": total})
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Create(&project)
	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var project models.Project
	if err := h.db.Preload("Client").Preload("Tasks").Preload("Labels").Preload("Members").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"}); return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var project models.Project
	if err := h.db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.ShouldBindJSON(&project)
	h.db.Save(&project)
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Project{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *ProjectHandler) GetTasks(c *gin.Context) {
	id, _ := getID(c)
	var tasks []models.Task
	h.db.Preload("AssignedTo").Where("project_id = ?", id).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func (h *ProjectHandler) GetTimeline(c *gin.Context) {
	id, _ := getID(c)
	var tasks []models.Task
	h.db.Where("project_id = ?", id).Order("updated_at desc").Limit(20).Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// ─── TASK ────────────────────────────────────────────

type TaskHandler struct{ db *gorm.DB }

func NewTaskHandler(db *gorm.DB) *TaskHandler { return &TaskHandler{db: db} }

func (h *TaskHandler) List(c *gin.Context) {
	var q PaginationQuery
	c.ShouldBindQuery(&q)
	var tasks []models.Task
	var total int64
	query := h.db.Model(&models.Task{}).Preload("AssignedTo").Preload("Project").Preload("Labels")
	if q.Q != "" {
		query = query.Where("title ILIKE ?", "%"+q.Q+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if assignedTo := c.Query("assigned_to_id"); assignedTo != "" {
		query = query.Where("assigned_to_id = ?", assignedTo)
	}
	query.Count(&total)
	query.Scopes(paginate(q)).Order("id desc").Find(&tasks)
	c.JSON(http.StatusOK, gin.H{"data": tasks, "total": total})
}

func (h *TaskHandler) Create(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Create(&task)
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var task models.Task
	if err := h.db.Preload("AssignedTo").Preload("Project").Preload("Collaborators").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"}); return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.ShouldBindJSON(&task)
	h.db.Save(&task)
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, _ := getID(c)
	var req struct{ Status string `json:"status"` }
	c.ShouldBindJSON(&req)
	h.db.Model(&models.Task{}).Where("id = ?", id).Update("status", req.Status)
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Task{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ─── LEAD ────────────────────────────────────────────

type LeadHandler struct{ db *gorm.DB }

func NewLeadHandler(db *gorm.DB) *LeadHandler { return &LeadHandler{db: db} }

func (h *LeadHandler) List(c *gin.Context) {
	var q PaginationQuery
	c.ShouldBindQuery(&q)
	var leads []models.Lead
	var total int64
	query := h.db.Model(&models.Lead{}).Preload("Owner").Preload("Labels")
	if q.Q != "" {
		query = query.Where("name ILIKE ?", "%"+q.Q+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	query.Scopes(paginate(q)).Order("id desc").Find(&leads)
	c.JSON(http.StatusOK, gin.H{"data": leads, "total": total})
}

func (h *LeadHandler) Create(c *gin.Context) {
	var lead models.Lead
	if err := c.ShouldBindJSON(&lead); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	lead.OwnerID = getUserID(c)
	h.db.Create(&lead)
	c.JSON(http.StatusCreated, lead)
}

func (h *LeadHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var lead models.Lead
	if err := h.db.Preload("Owner").First(&lead, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.JSON(http.StatusOK, lead)
}

func (h *LeadHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var lead models.Lead
	if err := h.db.First(&lead, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.ShouldBindJSON(&lead)
	h.db.Save(&lead)
	c.JSON(http.StatusOK, lead)
}

func (h *LeadHandler) UpdateStatus(c *gin.Context) {
	id, _ := getID(c)
	var req struct{ Status string `json:"status"` }
	c.ShouldBindJSON(&req)
	h.db.Model(&models.Lead{}).Where("id = ?", id).Update("status", req.Status)
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func (h *LeadHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Lead{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *LeadHandler) ConvertToClient(c *gin.Context) {
	id, _ := getID(c)
	var lead models.Lead
	if err := h.db.First(&lead, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lead not found"}); return
	}
	client := models.Client{
		Name:    lead.Name,
		Email:   lead.Email,
		Phone:   lead.Phone,
		OwnerID: getUserID(c),
	}
	if err := h.db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	h.db.Model(&lead).Update("status", "won")
	c.JSON(http.StatusCreated, client)
}

// ─── INVOICE ─────────────────────────────────────────

type InvoiceHandler struct{ db *gorm.DB }

func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler { return &InvoiceHandler{db: db} }

func (h *InvoiceHandler) List(c *gin.Context) {
	var q PaginationQuery
	c.ShouldBindQuery(&q)
	var invoices []models.Invoice
	var total int64
	query := h.db.Model(&models.Invoice{}).Preload("Client").Preload("Project").Preload("Labels")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	query.Scopes(paginate(q)).Order("id desc").Find(&invoices)
	c.JSON(http.StatusOK, gin.H{"data": invoices, "total": total})
}

func (h *InvoiceHandler) Create(c *gin.Context) {
	var invoice models.Invoice
	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Create(&invoice)
	c.JSON(http.StatusCreated, invoice)
}

func (h *InvoiceHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var invoice models.Invoice
	if err := h.db.Preload("Client").Preload("Project").Preload("Items").Preload("Payments").First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var invoice models.Invoice
	if err := h.db.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.ShouldBindJSON(&invoice)
	h.db.Save(&invoice)
	c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Invoice{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *InvoiceHandler) AddPayment(c *gin.Context) {
	id, _ := getID(c)
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	payment.InvoiceID = id
	h.db.Create(&payment)
	h.recalcInvoice(id)
	c.JSON(http.StatusCreated, payment)
}

func (h *InvoiceHandler) AddItem(c *gin.Context) {
	id, _ := getID(c)
	var item models.InvoiceItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	item.InvoiceID = id
	item.Total = item.Quantity * item.UnitPrice
	h.db.Create(&item)
	h.recalcInvoice(id)
	c.JSON(http.StatusCreated, item)
}

func (h *InvoiceHandler) UpdateItem(c *gin.Context) {
	invoiceID, _ := getID(c)
	itemID, _ := strconv.ParseUint(c.Param("itemId"), 10, 64)
	var item models.InvoiceItem
	if err := h.db.Where("id = ? AND invoice_id = ?", itemID, invoiceID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.ShouldBindJSON(&item)
	item.Total = item.Quantity * item.UnitPrice
	h.db.Save(&item)
	h.recalcInvoice(invoiceID)
	c.JSON(http.StatusOK, item)
}

func (h *InvoiceHandler) DeleteItem(c *gin.Context) {
	invoiceID, _ := getID(c)
	itemID, _ := strconv.ParseUint(c.Param("itemId"), 10, 64)
	h.db.Where("id = ? AND invoice_id = ?", itemID, invoiceID).Delete(&models.InvoiceItem{})
	h.recalcInvoice(invoiceID)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *InvoiceHandler) DeletePayment(c *gin.Context) {
	invoiceID, _ := getID(c)
	paymentID, _ := strconv.ParseUint(c.Param("paymentId"), 10, 64)
	var payment models.Payment
	if err := h.db.Where("id = ? AND invoice_id = ?", paymentID, invoiceID).First(&payment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	h.db.Delete(&payment)
	h.recalcInvoice(invoiceID)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *InvoiceHandler) recalcInvoice(invoiceID uint) {
	var items []models.InvoiceItem
	h.db.Where("invoice_id = ?", invoiceID).Find(&items)
	var subtotal float64
	for _, it := range items {
		subtotal += it.Total
	}
	var paidAmount float64
	h.db.Model(&models.Payment{}).Where("invoice_id = ?", invoiceID).Select("COALESCE(SUM(amount), 0)").Scan(&paidAmount)
	var invoice models.Invoice
	h.db.First(&invoice, invoiceID)
	invoice.TotalAmount = subtotal + invoice.TaxAmount - invoice.DiscountAmount
	invoice.PaidAmount = paidAmount
	invoice.DueAmount = invoice.TotalAmount - invoice.PaidAmount
	if invoice.TotalAmount > 0 && invoice.PaidAmount >= invoice.TotalAmount {
		invoice.Status = "fully_paid"
	} else if invoice.PaidAmount > 0 {
		invoice.Status = "partially_paid"
	} else if invoice.Status != "draft" && !invoice.DueDate.IsZero() && invoice.DueDate.Time.Before(time.Now()) {
		invoice.Status = "overdue"
	}
	h.db.Save(&invoice)
}

func (h *InvoiceHandler) Summary(c *gin.Context) {
	var summary []struct {
		ClientName    string  `json:"client_name"`
		Count         int     `json:"count"`
		InvoiceTotal  float64 `json:"invoice_total"`
		PaymentReceived float64 `json:"payment_received"`
		Due           float64 `json:"due"`
	}
	h.db.Table("invoices i").
		Select("c.name as client_name, COUNT(i.id) as count, SUM(i.total_amount) as invoice_total, SUM(i.paid_amount) as payment_received, SUM(i.due_amount) as due").
		Joins("JOIN clients c ON c.id = i.client_id").
		Where("i.deleted_at IS NULL").
		Group("c.name").
		Scan(&summary)
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// ─── PAYMENT ─────────────────────────────────────────

type PaymentHandler struct{ db *gorm.DB }

func NewPaymentHandler(db *gorm.DB) *PaymentHandler { return &PaymentHandler{db: db} }

func (h *PaymentHandler) List(c *gin.Context) {
	var payments []models.Payment
	var total int64
	h.db.Model(&models.Payment{}).Count(&total)
	h.db.Preload("Invoice").Order("payment_date desc").Find(&payments)
	c.JSON(http.StatusOK, gin.H{"data": payments, "total": total})
}

func (h *PaymentHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var payment models.Payment
	if err := h.db.Preload("Invoice").First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.JSON(http.StatusOK, payment)
}

// ─── CONTRACT ────────────────────────────────────────

type ContractHandler struct{ db *gorm.DB }

func NewContractHandler(db *gorm.DB) *ContractHandler { return &ContractHandler{db: db} }

func (h *ContractHandler) List(c *gin.Context) {
	var contracts []models.Contract
	var total int64
	h.db.Model(&models.Contract{}).Count(&total)
	h.db.Preload("Client").Preload("Project").Find(&contracts)
	c.JSON(http.StatusOK, gin.H{"data": contracts, "total": total})
}

func (h *ContractHandler) Create(c *gin.Context) {
	var contract models.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Create(&contract)
	c.JSON(http.StatusCreated, contract)
}

func (h *ContractHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var contract models.Contract
	if err := h.db.Preload("Client").Preload("Project").First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"}); return
	}
	c.JSON(http.StatusOK, contract)
}

func (h *ContractHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var contract models.Contract
	h.db.First(&contract, id)
	c.ShouldBindJSON(&contract)
	h.db.Save(&contract)
	c.JSON(http.StatusOK, contract)
}

func (h *ContractHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Contract{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ─── ITEM ────────────────────────────────────────────

type ItemHandler struct{ db *gorm.DB }

func NewItemHandler(db *gorm.DB) *ItemHandler { return &ItemHandler{db: db} }

func (h *ItemHandler) List(c *gin.Context) {
	var items []models.Item
	var total int64
	h.db.Model(&models.Item{}).Count(&total)
	h.db.Find(&items)
	c.JSON(http.StatusOK, gin.H{"data": items, "total": total})
}

func (h *ItemHandler) Create(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var item models.Item
	h.db.First(&item, id)
	c.ShouldBindJSON(&item)
	h.db.Save(&item)
	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Item{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ─── ORDER ───────────────────────────────────────────

type OrderHandler struct{ db *gorm.DB }

func NewOrderHandler(db *gorm.DB) *OrderHandler { return &OrderHandler{db: db} }

func (h *OrderHandler) List(c *gin.Context) {
	var orders []models.Order
	var total int64
	h.db.Model(&models.Order{}).Count(&total)
	h.db.Preload("Client").Find(&orders)
	c.JSON(http.StatusOK, gin.H{"data": orders, "total": total})
}

func (h *OrderHandler) Create(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Create(&order)
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var order models.Order
	h.db.Preload("Client").First(&order, id)
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var order models.Order
	h.db.First(&order, id)
	c.ShouldBindJSON(&order)
	h.db.Save(&order)
	c.JSON(http.StatusOK, order)
}

// ─── EVENT ───────────────────────────────────────────

type EventHandler struct{ db *gorm.DB }

func NewEventHandler(db *gorm.DB) *EventHandler { return &EventHandler{db: db} }

func (h *EventHandler) List(c *gin.Context) {
	var events []models.Event
	query := h.db.Preload("Labels")
	if month := c.Query("month"); month != "" {
		query = query.Where("EXTRACT(MONTH FROM start_date) = ?", month)
	}
	if year := c.Query("year"); year != "" {
		query = query.Where("EXTRACT(YEAR FROM start_date) = ?", year)
	}
	query.Find(&events)
	c.JSON(http.StatusOK, gin.H{"data": events})
}

func (h *EventHandler) Create(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	event.CreatedByID = getUserID(c)
	h.db.Create(&event)
	c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) Get(c *gin.Context) {
	id, _ := getID(c)
	var event models.Event
	h.db.Preload("Labels").First(&event, id)
	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var event models.Event
	h.db.First(&event, id)
	c.ShouldBindJSON(&event)
	h.db.Save(&event)
	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Event{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ─── NOTE ────────────────────────────────────────────

type NoteHandler struct{ db *gorm.DB }

func NewNoteHandler(db *gorm.DB) *NoteHandler { return &NoteHandler{db: db} }

func (h *NoteHandler) List(c *gin.Context) {
	userID := getUserID(c)
	var notes []models.Note
	h.db.Where("user_id = ?", userID).Find(&notes)
	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func (h *NoteHandler) Create(c *gin.Context) {
	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	note.UserID = getUserID(c)
	h.db.Create(&note)
	c.JSON(http.StatusCreated, note)
}

func (h *NoteHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var note models.Note
	h.db.First(&note, id)
	c.ShouldBindJSON(&note)
	h.db.Save(&note)
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Note{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ─── EXPENSE ─────────────────────────────────────────

type ExpenseHandler struct{ db *gorm.DB }

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler { return &ExpenseHandler{db: db} }

func (h *ExpenseHandler) List(c *gin.Context) {
	userID := getUserID(c)
	var expenses []models.Expense
	var total int64
	recurring := c.Query("recurring") == "true"
	h.db.Model(&models.Expense{}).Where("user_id = ? AND is_recurring = ?", userID, recurring).Count(&total)
	h.db.Where("user_id = ? AND is_recurring = ?", userID, recurring).Order("date desc").Find(&expenses)
	c.JSON(http.StatusOK, gin.H{"data": expenses, "total": total})
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	expense.UserID = getUserID(c)
	expense.Total = expense.Amount + expense.Tax + expense.SecondTax
	h.db.Create(&expense)
	c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	id, _ := getID(c)
	var expense models.Expense
	h.db.First(&expense, id)
	c.ShouldBindJSON(&expense)
	expense.Total = expense.Amount + expense.Tax + expense.SecondTax
	h.db.Save(&expense)
	c.JSON(http.StatusOK, expense)
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Expense{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ─── TEAM ────────────────────────────────────────────

type TeamHandler struct{ db *gorm.DB }

func NewTeamHandler(db *gorm.DB) *TeamHandler { return &TeamHandler{db: db} }

func (h *TeamHandler) ListMembers(c *gin.Context) {
	var members []models.User
	var total int64
	active := c.Query("inactive") != "true"
	h.db.Model(&models.User{}).Where("is_active = ?", active).Count(&total)
	h.db.Where("is_active = ?", active).Order("name asc").Find(&members)
	c.JSON(http.StatusOK, gin.H{"data": members, "total": total})
}

func (h *TeamHandler) GetMember(c *gin.Context) {
	id, _ := getID(c)
	var member models.User
	h.db.First(&member, id)
	c.JSON(http.StatusOK, member)
}

func (h *TeamHandler) UpdateMember(c *gin.Context) {
	id, _ := getID(c)
	var member models.User
	h.db.First(&member, id)
	c.ShouldBindJSON(&member)
	h.db.Save(&member)
	c.JSON(http.StatusOK, member)
}

func (h *TeamHandler) DeleteMember(c *gin.Context) {
	id, _ := getID(c)
	var member models.User
	if err := h.db.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}); return
	}
	h.db.Model(&member).Update("is_active", false)
	c.JSON(http.StatusOK, gin.H{"message": "User deactivated"})
}

func (h *TeamHandler) ResetPassword(c *gin.Context) {
	id, _ := getID(c)
	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"}); return
	}
	if err := h.db.Model(&models.User{}).Where("id = ?", id).Update("password", string(hashed)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"}); return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (h *TeamHandler) ListTimeCards(c *gin.Context) {
	var cards []models.TimeCard
	h.db.Preload("User").Order("in_date desc").Find(&cards)
	c.JSON(http.StatusOK, gin.H{"data": cards})
}

func (h *TeamHandler) ClockIn(c *gin.Context) {
	userID := getUserID(c)
	var existing models.TimeCard
	// Check if already clocked in
	if err := h.db.Where("user_id = ? AND out_time IS NULL", userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already clocked in"}); return
	}
	now := time.Now()
	card := models.TimeCard{UserID: userID, InTime: now, InDate: now}
	h.db.Create(&card)
	h.db.Model(&models.User{}).Where("id = ?", userID).Update("clocked_in", true)
	c.JSON(http.StatusCreated, card)
}

func (h *TeamHandler) ClockOut(c *gin.Context) {
	userID := getUserID(c)
	var card models.TimeCard
	if err := h.db.Where("user_id = ? AND out_time IS NULL", userID).First(&card).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active clock-in found"}); return
	}
	now := time.Now()
	duration := now.Sub(card.InTime).Hours()
	h.db.Model(&card).Updates(map[string]interface{}{
		"out_time": now,
		"out_date": now,
		"duration": duration,
	})
	h.db.Model(&models.User{}).Where("id = ?", userID).Update("clocked_in", false)
	c.JSON(http.StatusOK, gin.H{"message": "Clocked out"})
}

func (h *TeamHandler) ListLeaves(c *gin.Context) {
	userID := getUserID(c)
	var leaves []models.Leave
	h.db.Where("user_id = ?", userID).Order("start_date desc").Find(&leaves)
	c.JSON(http.StatusOK, gin.H{"data": leaves})
}

func (h *TeamHandler) ApplyLeave(c *gin.Context) {
	var leave models.Leave
	if err := c.ShouldBindJSON(&leave); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	leave.UserID = getUserID(c)
	leave.Status = "pending"
	h.db.Create(&leave)
	c.JSON(http.StatusCreated, leave)
}

func (h *TeamHandler) UpdateLeaveStatus(c *gin.Context) {
	id, _ := getID(c)
	var req struct{ Status string `json:"status"` }
	c.ShouldBindJSON(&req)
	h.db.Model(&models.Leave{}).Where("id = ?", id).Update("status", req.Status)
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

func (h *TeamHandler) ListAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	h.db.Preload("CreatedBy").Order("created_at desc").Find(&announcements)
	c.JSON(http.StatusOK, gin.H{"data": announcements})
}

func (h *TeamHandler) CreateAnnouncement(c *gin.Context) {
	var ann models.Announcement
	if err := c.ShouldBindJSON(&ann); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	ann.CreatedByID = getUserID(c)
	h.db.Create(&ann)
	c.JSON(http.StatusCreated, ann)
}

// ─── FILE ────────────────────────────────────────────

type FileHandler struct{ db *gorm.DB }

func NewFileHandler(db *gorm.DB) *FileHandler { return &FileHandler{db: db} }

func (h *FileHandler) List(c *gin.Context) {
	userID := getUserID(c)
	var files []models.File
	folderID := c.Query("folder_id")
	query := h.db.Where("owner_id = ?", userID)
	if folderID == "" {
		query = query.Where("folder_id IS NULL")
	} else {
		query = query.Where("folder_id = ?", folderID)
	}
	query.Find(&files)
	c.JSON(http.StatusOK, gin.H{"data": files})
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"}); return
	}
	defer file.Close()
	// TODO: implement actual file storage (S3, local, etc.)
	f := models.File{
		Name:     header.Filename,
		Size:     header.Size,
		MimeType: header.Header.Get("Content-Type"),
		OwnerID:  getUserID(c),
	}
	h.db.Create(&f)
	c.JSON(http.StatusCreated, f)
}

func (h *FileHandler) CreateFolder(c *gin.Context) {
	var folder models.File
	if err := c.ShouldBindJSON(&folder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	folder.IsFolder = true
	folder.OwnerID = getUserID(c)
	h.db.Create(&folder)
	c.JSON(http.StatusCreated, folder)
}

func (h *FileHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.File{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *FileHandler) ToggleFavorite(c *gin.Context) {
	id, _ := getID(c)
	var f models.File
	h.db.First(&f, id)
	h.db.Model(&f).Update("is_favorite", !f.IsFavorite)
	c.JSON(http.StatusOK, gin.H{"is_favorite": !f.IsFavorite})
}

// ─── TODO ────────────────────────────────────────────

type TodoHandler struct{ db *gorm.DB }

func NewTodoHandler(db *gorm.DB) *TodoHandler { return &TodoHandler{db: db} }

func (h *TodoHandler) List(c *gin.Context) {
	userID := getUserID(c)
	done := c.Query("done") == "true"
	var todos []models.Todo
	h.db.Where("user_id = ? AND done = ?", userID, done).Order("created_at desc").Find(&todos)
	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func (h *TodoHandler) Create(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	todo.UserID = getUserID(c)
	h.db.Create(&todo)
	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) MarkDone(c *gin.Context) {
	id, _ := getID(c)
	h.db.Model(&models.Todo{}).Where("id = ?", id).Updates(map[string]interface{}{"done": true, "done_at": gorm.Expr("NOW()")})
	c.JSON(http.StatusOK, gin.H{"message": "Marked as done"})
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Todo{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// ─── REPORT ──────────────────────────────────────────

type ReportHandler struct{ db *gorm.DB }

func NewReportHandler(db *gorm.DB) *ReportHandler { return &ReportHandler{db: db} }

func (h *ReportHandler) InvoicesSummary(c *gin.Context) {
	year := c.Query("year")
	var summary []struct {
		ClientName      string  `json:"client_name"`
		Count           int     `json:"count"`
		InvoiceTotal    float64 `json:"invoice_total"`
		TaxAmount       float64 `json:"tax_amount"`
		PaymentReceived float64 `json:"payment_received"`
		Due             float64 `json:"due"`
	}
	query := h.db.Table("invoices i").
		Select("c.name as client_name, COUNT(i.id) as count, SUM(i.total_amount) as invoice_total, SUM(i.tax_amount) as tax_amount, SUM(i.paid_amount) as payment_received, SUM(i.due_amount) as due").
		Joins("JOIN clients c ON c.id = i.client_id").
		Where("i.deleted_at IS NULL")
	if year != "" {
		query = query.Where("EXTRACT(YEAR FROM i.bill_date) = ?", year)
	}
	query.Group("c.name").Order("invoice_total desc").Scan(&summary)
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *ReportHandler) ProjectsSummary(c *gin.Context) {
	var summary struct {
		Open      int64 `json:"open"`
		Completed int64 `json:"completed"`
		Hold      int64 `json:"hold"`
		Cancelled int64 `json:"cancelled"`
	}
	h.db.Model(&models.Project{}).Where("status = 'open'").Count(&summary.Open)
	h.db.Model(&models.Project{}).Where("status = 'completed'").Count(&summary.Completed)
	h.db.Model(&models.Project{}).Where("status = 'hold'").Count(&summary.Hold)
	h.db.Model(&models.Project{}).Where("status = 'cancelled'").Count(&summary.Cancelled)
	c.JSON(http.StatusOK, summary)
}

func (h *ReportHandler) LeadsSummary(c *gin.Context) {
	var summary []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	h.db.Model(&models.Lead{}).Select("status, COUNT(*) as count").Group("status").Scan(&summary)
	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *ReportHandler) ExpensesSummary(c *gin.Context) {
	var total float64
	h.db.Model(&models.Expense{}).Select("COALESCE(SUM(total), 0)").Scan(&total)
	c.JSON(http.StatusOK, gin.H{"total": total})
}

// ─── LABEL ───────────────────────────────────────────

type LabelHandler struct{ db *gorm.DB }

func NewLabelHandler(db *gorm.DB) *LabelHandler { return &LabelHandler{db: db} }

func (h *LabelHandler) List(c *gin.Context) {
	var labels []models.Label
	h.db.Find(&labels)
	c.JSON(http.StatusOK, gin.H{"data": labels})
}

func (h *LabelHandler) Create(c *gin.Context) {
	var label models.Label
	if err := c.ShouldBindJSON(&label); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return
	}
	h.db.Create(&label)
	c.JSON(http.StatusCreated, label)
}

func (h *LabelHandler) Delete(c *gin.Context) {
	id, _ := getID(c)
	h.db.Delete(&models.Label{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
