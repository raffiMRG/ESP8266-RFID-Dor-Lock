#pragma once

#include <Arduino.h>

// ---- RFID RC522 (SPI) ----
// SCK=D5(GPIO14), MOSI=D7(GPIO13), MISO=D6(GPIO12) are the fixed HW SPI pins,
// wire them directly to the RC522 module.
#define RFID_SS_PIN D4  // GPIO2  -> RC522 SDA/SS
#define RFID_RST_PIN D3 // GPIO0  -> RC522 RST

// ---- I2C LCD 16x2 ----
// SDA=D2(GPIO4), SCL=D1(GPIO5) are the default HW I2C pins on the ESP8266 Arduino core.
#define LCD_I2C_ADDR 0x27
#define LCD_COLS 16
#define LCD_ROWS 2

// ---- Relay (solenoid door lock via SRD-05VDC-SL-C module) ----
#define RELAY_PIN D0          // GPIO16 -> relay IN
#define RELAY_ACTIVE_LOW true // most SRD-05VDC-SL-C boards energize the relay on LOW
#define UNLOCK_DURATION_MS 5000UL

// ---- Buzzer (active buzzer module, on/off only) ----
#define BUZZER_PIN D8           // GPIO15 -> buzzer I/O
#define BUZZER_ACTIVE_LOW false // most active-buzzer modules sound on HIGH
#define BEEP_SHORT_MS 100UL
#define BEEP_GAP_MS 100UL

// ---- WiFi ----
#define WIFI_SSID "UB"
#define WIFI_PASSWORD "123123123"
#define WIFI_CONNECT_TIMEOUT_MS 15000UL

// ---- Backend API ----
// Otorisasi & pencatatan log akses dilakukan di backend (lihat /backend),
// bukan lagi hardcoded di firmware. Kelola kartu lewat dashboard frontend.
#define API_BASE_URL "http://192.168.1.123:8765" // IP server yang menjalankan docker compose (BACKEND_PORT di .env)
#define API_ACCESS_PATH "/api/v1/access"
#define API_UNLOCK_POLL_PATH "/api/v1/unlock/poll"
#define HTTP_TIMEOUT_MS 5000
#define REMOTE_UNLOCK_POLL_MS 1500UL // jeda antar polling status buka pintu jarak jauh
