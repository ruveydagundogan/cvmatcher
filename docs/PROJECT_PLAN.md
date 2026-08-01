# CV Matcher v2 — Proje Planı

> Tarih: 2026-07-30 · Durum: Onaylandı · Sahip: PM

## Kurgu

Uygulama iki persona tarafından kullanılır:

| Persona | Rol | Dünya |
|---|---|---|
| Bireysel iş arayan | `user` | Kendi CV'leri, JD'leri, match skorları ve CV Asistanı chat |
| İK uzmanı (işveren) | `hr` | Kendi eklediği adaylar, JD'ler, pipeline, en iyi aday eşleştirme |
| Teknik yönetici | `admin` | Model/adapter, prompt, ayar ve kullanıcı yönetimi |

## Tasarım Kararları

- İK kurgusu: **tek şirket, tek İK hesabı**
- Aday CV'leri: İK **kendisi ekler** (başvuru sistemi yok) — aday verileri bireysel kullanıcılardan tamamen ayrı ve private
- Chat: ayrı **cv-coach** fine-tune modeli
- Sıra: P0 → P1 (Chat) → P2 (İK paneli) → P3 (iyileştirmeler)

## Faz 0 — Temel Sağlamlaştırma

- [x] F0-3: JD analiz modeli kararı — `cv-parser` korundu (ikisi de eşit başarı; cv-parser yapılandırılmış çıkarım için eğitilmiş)
- [ ] F0-1: Admin repo → Postgres (restart'ta veri kalıcılığı)
- [ ] F0-2: Rol bazlı frontend routing (`/admin` admin'e, `/ik` hr'ya)
- [ ] F0-4: Kayıtta rol seçimi (`user` / `hr`)

## Faz 1 — CV Asistanı Chat (bireysel)

### 1a. cv-coach modeli
- `train_lora.py`'ye `--mode cv-coach` eklenir
- Eğitim verisi: CV geliştirme konuşmaları (Q&A çiftleri)
- Konular: eksik skill tespiti, deneyim yazımı, özet iyileştirme, JD'ye göre uyarlama, mülakat hazırlığı
- Pipeline: safetensors → GGUF → Ollama `cv-coach` modeli

### 1b. Backend
- Migration `008`: `conversations` + `messages`
- API: `GET/POST /api/v1/chat/conversations`, `DELETE /{id}`, `POST /{id}/messages` (SSE stream)
- Sistem promptu kullanıcının parse edilmiş CV verisini içerir; son N mesaj bağlam
- `ChatCompletionWithModel` stream varyantı

### 1c. Frontend
- `/dashboard/coach`: solda konuşma geçmişi, sağda sohbet, SSE akışı
- Mevcut ölü kod (`PromptInput`, `useLLM`) burada kullanılır

## Faz 2 — İK Paneli (hr)

### 2a. Veritabanı — Migration `009`
- `candidates(id, user_id(hr), title, content, parsed_*, status, notes, created_at)`
- `status`: `new → screening → interview → offer → hired / rejected`

### 2b. Backend
- `POST /api/v1/ik/candidates` (CV ekle + otomatik parse)
- `GET /api/v1/ik/candidates` (filtre: status, skill, skor)
- `PATCH /api/v1/ik/candidates/{id}` (durum + not)
- `POST /api/v1/ik/jds/{id}/find-candidates` (en iyi aday sıralaması)
- `GET /api/v1/ik/stats` (pipeline dağılımı)

### 2c. Frontend — `/ik` alanı
- Aday listesi, aday ekle, detay, pipeline, "JD'ye en iyi aday"

## Faz 3 — İyileştirmeler

- CV şablonları & sürüm geçmişi
- PDF CV yükleme
- E-posta bildirimleri
- Admin'den kullanıcı yönetimi

## Teknik Borç Notları

1. Admin repo in-memory → Postgres (F0-1)
2. Frontend admin layout rol kontrolü yapmıyor (F0-2)
3. JD analiz `cv-parser` kullanıyor (F0-3, tamamlandı)
4. `llmscoring` legacy modülü aktif ama frontend'de kullanılmıyor
5. README hâlâ Gemma-2B referanslı

## İş Yükü Tahmini

| Faz | Backend | Frontend | Eğitim | Toplam |
|---|---|---|---|---|
| F0 | ~2s | ~2s | — | ~4s |
| F1 | ~4s | ~4s | ~2s | ~10s |
| F2 | ~5s | ~5s | — | ~10s |
| F3 | ~4s | ~4s | — | ~8s |
