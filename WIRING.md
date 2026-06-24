# Wiring — RFID Door Lock (ESP8266 NodeMCU Amica)

Board: **NodeMCU ESP8266 (Amica, ESP-12E)** — PlatformIO id `nodemcuv2`

> ⚠️ Semua pin GPIO ESP8266 beroperasi di level **3.3V**. Jangan hubungkan pin sinyal apa pun ke 5V.

## 1. RFID RC522 (SPI, 3.3V)

| Pin RC522 | Pin NodeMCU | GPIO     |
|-----------|--------------|----------|
| SDA (SS)  | D4           | GPIO2    |
| SCK       | D5           | GPIO14   |
| MOSI      | D7           | GPIO13   |
| MISO      | D6           | GPIO12   |
| IRQ       | *(tidak dipakai)* | — |
| GND       | GND          | —        |
| RST       | D3           | GPIO0    |
| 3.3V      | 3V3          | —        |

> ⚠️ RC522 hanya tahan **3.3V**, jangan hubungkan ke pin 5V/VU/VIN — bisa merusak modul.

## 2. LCD 16x2 I2C (backpack PCF8574, alamat default `0x27`)

| Pin LCD I2C | Pin NodeMCU | GPIO   |
|-------------|--------------|--------|
| SDA         | D2           | GPIO4  |
| SCL         | D1           | GPIO5  |
| VCC         | VU (5V)      | —      |
| GND         | GND          | —      |

> Backpack PCF8574 umumnya butuh 5V untuk backlight optimal — NodeMCU pin `VU`/`VIN` menyediakan 5V langsung dari USB. Jalur I2C (SDA/SCL) tetap aman karena PCF8574 mentoleransi logic 3.3V dari ESP8266.

## 3. Relay 1-channel (SRD-05VDC-SL-C) → Solenoid Door Lock

| Pin Relay | Pin NodeMCU | GPIO   |
|-----------|--------------|--------|
| IN        | D0           | GPIO16 |
| VCC       | VU (5V)      | —      |
| GND       | GND          | —      |

**Sisi output relay (tegangan tinggi/solenoid, terpisah dari NodeMCU):**

| Pin Relay | Terhubung ke |
|-----------|--------------|
| COM       | (+) Power Supply solenoid (umumnya 12V DC) |
| NO        | (+) Solenoid lock |
| —         | (−) Solenoid lock → (−) Power Supply (kembali ke GND power supply, **bukan** GND NodeMCU) |

> Gunakan jalur **NO (Normally Open)** supaya pintu tetap terkunci (default) saat ESP8266 mati/reset, dan hanya terbuka saat relay aktif.
> Power supply solenoid (12V) harus **terpisah** dari suplai 5V NodeMCU — relay module mengisolasi sisi kontrol (3.3V/5V logic) dari sisi beban (12V) lewat optocoupler/coil, jangan disatukan GND-nya kecuali memang didesain common ground oleh modul relay yang dipakai.

## 4. Buzzer (active buzzer module, on/off)

| Pin Buzzer | Pin NodeMCU | GPIO   |
|------------|--------------|--------|
| I/O / S    | D8           | GPIO15 |
| VCC        | 3V3 atau VU (5V, sesuai modul) | — |
| GND        | GND          | —      |

> D8 (GPIO15) harus **LOW** saat boot agar ESP8266 boot normal. Modul active-buzzer 3-pin (VCC/GND/I-O dengan transistor) umumnya menarik pin ini LOW secara default sehingga aman dipakai di sini — hindari menambahkan pull-up eksternal pada pin ini.

## Diagram ringkas

```
                ┌──────────────────────┐
   RC522        │                      │        LCD I2C 16x2
 ┌─────────┐    │   NodeMCU ESP8266    │      ┌──────────────┐
 │ SDA ─────────┤ D4               D2  ├──────│ SDA          │
 │ SCK ─────────┤ D5               D1  ├──────│ SCL          │
 │ MOSI─────────┤ D7               VU  ├──────│ VCC          │
 │ MISO─────────┤ D6               GND ├──────│ GND          │
 │ RST ─────────┤ D3                   │      └──────────────┘
 │ GND ─────────┤ GND                  │
 │ 3.3V─────────┤ 3V3              D0  ├───┐
 └─────────┘    └──────────────────────┘   │
                                            │   Relay SRD-05VDC-SL-C
                                       ┌────┴───────┐
                                       │ IN         │      Solenoid Lock
                                       │ VCC ← VU   │     ┌───────────┐
                                       │ GND ← GND  │     │           │
                                       │            │     │           │
                                       │ COM ───────┼─────┤ (+) PSU 12V
                                       │ NO  ───────┼─────┤ (+) Solenoid
                                       └────────────┘     │ (−) Solenoid ── (−) PSU 12V
                                                           └───────────┘

   NodeMCU D8 ── I/O Buzzer (active buzzer module)
   NodeMCU 3V3/VU ── VCC Buzzer
   NodeMCU GND ── GND Buzzer
```

## Referensi pin NodeMCU (Amica)

| Label board | GPIO    | Catatan boot                          |
|-------------|---------|----------------------------------------|
| D0          | GPIO16  | Tidak punya pull-up internal           |
| D1          | GPIO5   | I2C SCL                                 |
| D2          | GPIO4   | I2C SDA                                 |
| D3          | GPIO0   | Harus HIGH saat boot (flash mode select) — RST RC522 idle HIGH, aman |
| D4          | GPIO2   | Harus HIGH saat boot — LED onboard, SS RC522 idle HIGH, aman |
| D5          | GPIO14  | HW SPI SCK                              |
| D6          | GPIO12  | HW SPI MISO                             |
| D7          | GPIO13  | HW SPI MOSI                             |
| D8          | GPIO15  | Harus LOW saat boot — dipakai untuk buzzer (aman karena modul active-buzzer default LOW) |

Konfigurasi pin lengkap ada di [`include/config.h`](include/config.h).
