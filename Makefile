# ==========================================
# CodeMeet Go Library Setup Configuration
# ==========================================

APP_NAME := codemeet-bot
VERSION  := 1.0.1
AUTHOR   := Abolfazl Zarei
GITHUB   := github.com/AbolfazlZarei-dev/codemeet-go

# Go related variables
GOCMD   := go
GOBUILD := $(GOCMD) build
GORUN   := $(GOCMD) run
GOTEST  := $(GOCMD) test
GOMOD   := $(GOCMD) mod
GOFMT   := $(GOCMD) fmt
GOVET   := $(GOCMD) vet

# ==========================================
# Targets (Commands)
# ==========================================

.PHONY: all info deps run build test clean tidy fmt vet check help

# نمایش اطلاعات پروژه (پیش‌فرض هنگام اجرای make)
all: info

# ==========================================
# نمایش اطلاعات کتابخانه
# ==========================================

info:
	@echo "--------------------------------------------------"
	@echo " CodeMeet Go Bot"
	@echo "--------------------------------------------------"
	@echo " Name        : $(APP_NAME)"
	@echo " Version     : v$(VERSION)"
	@echo " Author      : $(AUTHOR)"
	@echo " GitHub      : $(GITHUB)"
	@echo " Go Version  : 1.20+"
	@echo "--------------------------------------------------"

# ==========================================
# نصب / دریافت وابستگی‌ها
# ==========================================

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# ==========================================
# مرتب‌سازی وابستگی‌ها
# ==========================================

tidy:
	@echo "Tidying up modules..."
	$(GOMOD) tidy

# ==========================================
# اجرای ربات نمونه
# ==========================================

run:
	@echo "Running the example bot..."
	$(GORUN) ./examples

# ==========================================
# ساخت فایل اجرایی ربات
# ==========================================

build:
	@echo "Building the bot binary..."
	$(GOBUILD) -o ./$(APP_NAME) ./examples
	@echo "Build successful!"
	@echo "Output: ./$(APP_NAME)"

# ==========================================
# اجرای تست‌ها
# ==========================================

test:
	@echo "Running tests..."
	$(GOTEST) ./...

# ==========================================
# بررسی کامل پروژه
# ==========================================

check:
	@echo "Running project checks..."
	$(GOTEST) ./...
	$(GOVET) ./...
	@echo "All checks completed successfully!"

# ==========================================
# پاک کردن فایل‌های ساخته‌شده
# ==========================================

clean:
	@echo "Cleaning up..."
	$(GOCMD) clean
	rm -f $(APP_NAME)
	@echo "Clean completed!"

# ==========================================
# فرمت کردن کدهای Go
# ==========================================

fmt:
	@echo "Formatting Go files..."
	$(GOFMT) ./...
	@echo "Formatting completed!"

# ==========================================
# بررسی خطاهای استاتیک
# ==========================================

vet:
	@echo "Vetting Go files..."
	$(GOVET) ./...
	@echo "Vet completed successfully!"

# ==========================================
# راهنمای دستورات
# ==========================================

help:
	@echo "Available commands:"
	@echo ""
	@echo "  make info   - Show project information"
	@echo "  make deps   - Download dependencies"
	@echo "  make tidy   - Clean up go.mod and go.sum"
	@echo "  make run    - Run the example bot"
	@echo "  make build  - Build the example bot"
	@echo "  make test   - Run all tests"
	@echo "  make check  - Run tests and go vet"
	@echo "  make clean  - Remove built binaries"
	@echo "  make fmt    - Format Go code"
	@echo "  make vet    - Run go vet"