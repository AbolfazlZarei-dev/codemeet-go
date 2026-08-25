# ==========================================
# CodeMeet Go Library Setup Configuration
# ==========================================
APP_NAME := codemeet-bot
VERSION  := 1.0.0
AUTHOR   := Abolfazl Zarei
GITHUB   := github.com/AbolfazlZarei-dev/codemeet

# Go related variables
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod

# ==========================================
# Targets (Commands)
# ==========================================

.PHONY: all info deps run build clean tidy help

# نمایش اطلاعات پروژه (پیش‌فرض هنگام اجرای make)
all: info

# نمایش اطلاعات کتابخانه
info:
    @echo "--------------------------------------------------"
    @echo " CodeMeet Go Bot Setup"
    @echo "--------------------------------------------------"
    @echo " Name        : $(APP_NAME)"
    @echo " Version     : $(VERSION)"
    @echo " Author      : $(AUTHOR)"
    @echo " GitHub      : $(GITHUB)"
    @echo "--------------------------------------------------"

# نصب وابستگی‌ها
deps:
    @echo "Downloading dependencies..."
    $(GOMOD) download

# مرتب‌سازی وابستگی‌ها
tidy:
    @echo "Tidying up modules..."
    $(GOMOD) tidy

# اجرای مستقیم ربات نمونه
run:
    @echo "Running the bot..."
    cd examples && $(GORUN) .

# ساخت فایل اجرایی ربات (خروجی باینری)
build:
    @echo "Building the bot binary..."
    cd examples && $(GOBUILD) -o ../$(APP_NAME) .
    @echo "Build successful! Output: ./$(APP_NAME)"

# پاک کردن فایل‌های باینری اضافی
clean:
    @echo "Cleaning up..."
    $(GOCMD) clean
    rm -f $(APP_NAME)

# راهنمای دستورات
help:
    @echo "Available commands:"
    @echo "  make info   - Show project information"
    @echo "  make deps   - Download dependencies"
    @echo "  make tidy   - Clean up go.mod and go.sum"
    @echo "  make run    - Run the example bot directly"
    @echo "  make build  - Build the bot into an executable"
    @echo "  make clean  - Remove built binaries"