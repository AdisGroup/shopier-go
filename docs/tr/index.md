---
layout: home

hero:
  name: "Shopier Go SDK"
  text: "Shopier REST API için Üretime Hazır Go İstemcisi"
  tagline: "Sıfır dış bağımlılık. Eksiksiz API desteği. Go 1.23+ iteratörleri. Otomatik retry mekanizması."
  actions:
    - theme: brand
      text: Başlangıç Rehberi
      link: /tr/getting-started
    - theme: alt
      text: API Referansı
      link: /tr/resources/orders
    - theme: alt
      text: GitHub
      link: https://github.com/AdisGroup/shopier-go

features:
  - title: Sıfır Harici Bağımlılık
    details: Tamamen Go standart kütüphanesi (net/http, crypto/hmac, log/slog) ile geliştirildi. Güvenlik açığı veya paket şişkinliği barındırmaz.
  - title: Hata Toleranslı Mimari
    details: Üstel geri çekilme (exponential backoff), 60 saniyelik kayan pencere limit takibi ve Retry-After başlığı desteği.
  - title: Modern Go İteratörleri
    details: Klasik sayfalama yapısının yanında Go 1.23+ Range-over-func (client.Orders.All) ile doğrudan for döngüsü desteği.
  - title: Güvenli Webhook Doğrulama
    details: Sabit zamanlı HS256 HMAC imza kontrolü ve net/http, Gin, Chi, Fiber uyumlu hazır handler yapısı.
---
