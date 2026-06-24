Siap, ini perintah yang perlu Anda jalankan sendiri di terminal, dari folder project:

cd /home/mrg/code/iot/rfid-dor-lock

1. Compile saja (cek error tanpa upload):
   pio run

2. Compile + upload ke NodeMCU (pastikan board tersambung USB, terdeteksi di /dev/ttyUSB0):
   pio run -t upload

3. Lihat log serial setelah upload (baud 115200, keluar dengan Ctrl+C):
   pio device monitor

Atau gabung upload + langsung monitor:
pio run -t upload && pio device monitor

Kalau command pio tidak ditemukan

Saya sudah menambahkan ~/.platformio/penv/bin ke PATH lewat ~/.bashrc sebelumnya — kalau ini terminal baru, seharusnya sudah otomatis terbaca. Kalau belum:
source ~/.bashrc
atau buka terminal baru.

Project sudah pernah berhasil di-build & upload sebelumnya, jadi seharusnya langsung jalan dengan konfigurasi config.h Anda saat ini (WiFi Trojan.sh + backend 192.168.0.143:8080). Pastikan saja docker compose up -d di stack backend masih jalan supaya bisa diakses NodeMCU saat tes tap kartu.
