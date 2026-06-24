#include <Arduino.h>
#include <SPI.h>
#include <MFRC522.h>
#include <Wire.h>
#include <LiquidCrystal_I2C.h>
#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClient.h>
#include <ArduinoJson.h>

#include "config.h"

enum class AccessResult { Granted, Denied, ServerError };

MFRC522 rfid(RFID_SS_PIN, RFID_RST_PIN);
LiquidCrystal_I2C lcd(LCD_I2C_ADDR, LCD_COLS, LCD_ROWS);

void setLock(bool unlocked) {
  uint8_t activeLevel = RELAY_ACTIVE_LOW ? LOW : HIGH;
  uint8_t idleLevel = RELAY_ACTIVE_LOW ? HIGH : LOW;
  digitalWrite(RELAY_PIN, unlocked ? activeLevel : idleLevel);
}

void setBuzzer(bool on) {
  uint8_t activeLevel = BUZZER_ACTIVE_LOW ? LOW : HIGH;
  uint8_t idleLevel = BUZZER_ACTIVE_LOW ? HIGH : LOW;
  digitalWrite(BUZZER_PIN, on ? activeLevel : idleLevel);
}

void beep(uint16_t durationMs) {
  setBuzzer(true);
  delay(durationMs);
  setBuzzer(false);
}

void beepGrant() {
  beep(BEEP_SHORT_MS);
}

void beepDeny() {
  beep(BEEP_SHORT_MS);
  delay(BEEP_GAP_MS);
  beep(BEEP_SHORT_MS);
}

void showIdleScreen() {
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Tempelkan Kartu");
  lcd.setCursor(0, 1);
  lcd.print("Akses Pintu");
}

String formatUid(const MFRC522::Uid &uid) {
  String result;
  for (byte i = 0; i < uid.size; i++) {
    if (uid.uidByte[i] < 0x10) result += '0';
    result += String(uid.uidByte[i], HEX);
    if (i < uid.size - 1) result += ':';
  }
  result.toUpperCase();
  return result;
}

void connectWiFi() {
  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Menghubungkan");
  lcd.setCursor(0, 1);
  lcd.print("WiFi...");

  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  Serial.print(F("Menghubungkan ke WiFi"));
  unsigned long start = millis();
  while (WiFi.status() != WL_CONNECTED && millis() - start < WIFI_CONNECT_TIMEOUT_MS) {
    delay(300);
    Serial.print('.');
  }
  Serial.println();

  if (WiFi.status() == WL_CONNECTED) {
    Serial.print(F("WiFi terhubung, IP: "));
    Serial.println(WiFi.localIP());
  } else {
    Serial.println(F("WiFi gagal terhubung, lanjut tanpa koneksi"));
  }
}

AccessResult checkAccessOnServer(const String &uid, String &outLabel) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println(F("WiFi tidak terhubung"));
    return AccessResult::ServerError;
  }

  WiFiClient client;
  HTTPClient http;
  String url = String(API_BASE_URL) + API_ACCESS_PATH;

  if (!http.begin(client, url)) {
    Serial.println(F("Gagal memulai koneksi HTTP"));
    return AccessResult::ServerError;
  }

  http.addHeader("Content-Type", "application/json");
  http.setTimeout(HTTP_TIMEOUT_MS);

  String body = String("{\"uid\":\"") + uid + "\"}";
  int statusCode = http.POST(body);

  if (statusCode != 200) {
    Serial.print(F("Backend merespons status: "));
    Serial.println(statusCode);
    http.end();
    return AccessResult::ServerError;
  }

  String response = http.getString();
  http.end();

  JsonDocument doc;
  DeserializationError err = deserializeJson(doc, response);
  if (err) {
    Serial.print(F("Gagal parsing JSON: "));
    Serial.println(err.c_str());
    return AccessResult::ServerError;
  }

  bool granted = doc["granted"] | false;
  outLabel = doc["label"] | "";

  return granted ? AccessResult::Granted : AccessResult::Denied;
}

bool checkRemoteUnlock() {
  if (WiFi.status() != WL_CONNECTED) {
    return false;
  }

  WiFiClient client;
  HTTPClient http;
  String url = String(API_BASE_URL) + API_UNLOCK_POLL_PATH;

  if (!http.begin(client, url)) {
    return false;
  }
  http.setTimeout(HTTP_TIMEOUT_MS);

  int statusCode = http.GET();
  if (statusCode != 200) {
    http.end();
    return false;
  }

  String response = http.getString();
  http.end();

  JsonDocument doc;
  if (deserializeJson(doc, response)) {
    return false;
  }

  return doc["unlock"] | false;
}

void grantAccess(const char *label) {
  Serial.print(F("Akses diterima: "));
  Serial.println(label);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Akses Diterima");
  lcd.setCursor(0, 1);
  lcd.print(label);

  beepGrant();
  setLock(true);
  delay(UNLOCK_DURATION_MS);
  setLock(false);

  showIdleScreen();
}

void denyAccess(const char *reason) {
  Serial.print(F("Akses ditolak: "));
  Serial.println(reason);

  lcd.clear();
  lcd.setCursor(0, 0);
  lcd.print("Akses Ditolak");
  lcd.setCursor(0, 1);
  lcd.print(reason);
  beepDeny();
  delay(1500);

  showIdleScreen();
}

void setup() {
  Serial.begin(115200);

  pinMode(RELAY_PIN, OUTPUT);
  setLock(false);

  pinMode(BUZZER_PIN, OUTPUT);
  setBuzzer(false);

  SPI.begin();
  rfid.PCD_Init();
  rfid.PCD_DumpVersionToSerial();

  Wire.begin();
  lcd.init();
  lcd.backlight();

  connectWiFi();
  showIdleScreen();

  Serial.println(F("RFID Door Lock siap. Tempelkan kartu RFID..."));
}

void loop() {
  static unsigned long lastUnlockPoll = 0;
  if (millis() - lastUnlockPoll >= REMOTE_UNLOCK_POLL_MS) {
    lastUnlockPoll = millis();
    if (checkRemoteUnlock()) {
      grantAccess("Remote Unlock");
    }
  }

  if (!rfid.PICC_IsNewCardPresent() || !rfid.PICC_ReadCardSerial()) {
    return;
  }

  String uid = formatUid(rfid.uid);
  Serial.print(F("UID kartu: "));
  Serial.println(uid);

  String label;
  AccessResult result = checkAccessOnServer(uid, label);

  switch (result) {
    case AccessResult::Granted:
      grantAccess(label.c_str());
      break;
    case AccessResult::Denied:
      denyAccess("Kartu Tidak Dikenal");
      break;
    case AccessResult::ServerError:
      denyAccess("Server Error");
      break;
  }

  rfid.PICC_HaltA();
  rfid.PCD_StopCrypto1();
}
