import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Shopier Go SDK',
  description: 'Production-ready, zero-dependency Go SDK for the Shopier REST API',
  lastUpdated: true,
  cleanUrls: true,

  themeConfig: {
    logo: '/logo.svg',
    search: {
      provider: 'local',
      options: {
        locales: {
          tr: {
            translations: {
              button: {
                buttonText: 'Ara',
                buttonAriaLabel: 'Arama Yap'
              },
              modal: {
                noResultsText: 'Sonuç bulunamadı',
                resetButtonTitle: 'Aramayı sıfırla',
                footer: {
                  selectText: 'seç',
                  navigateText: 'gezin',
                  closeText: 'kapat'
                }
              }
            }
          }
        }
      }
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/AdisGroup/shopier-go' }
    ]
  },

  locales: {
    root: {
      label: 'English',
      lang: 'en',
      themeConfig: {
        nav: [
          { text: 'Guide', link: '/getting-started' },
          { text: 'API Reference', link: '/resources/orders' },
          { text: 'Webhooks', link: '/webhooks' },
          { text: 'Examples', link: 'https://github.com/AdisGroup/shopier-go/tree/main/examples' }
        ],
        sidebar: [
          {
            text: 'Getting Started',
            items: [
              { text: 'Overview & Installation', link: '/getting-started' },
              { text: 'Authentication (PAT & OAuth)', link: '/authentication' },
              { text: 'Pagination & Iterators', link: '/pagination-and-iterators' },
              { text: 'Error Handling & Retries', link: '/error-handling' },
              { text: 'Testing & Mocking', link: '/testing' }
            ]
          },
          {
            text: 'Webhooks',
            items: [
              { text: 'Signature Verification & Events', link: '/webhooks' }
            ]
          },
          {
            text: 'API Resources',
            items: [
              { text: 'Orders', link: '/resources/orders' },
              { text: 'Products', link: '/resources/products' },
              { text: 'Categories', link: '/resources/categories' },
              { text: 'Variations', link: '/resources/variations' },
              { text: 'Selections', link: '/resources/selections' },
              { text: 'Discounts', link: '/resources/discounts' },
              { text: 'Shippings', link: '/resources/shippings' },
              { text: 'Balance', link: '/resources/balance' },
              { text: 'Payouts', link: '/resources/payouts' },
              { text: 'Refunds', link: '/resources/refunds' },
              { text: 'Shop Settings', link: '/resources/shop' },
              { text: 'Webhooks Subscriptions', link: '/resources/webhooks' }
            ]
          }
        ]
      }
    },
    tr: {
      label: 'Türkçe',
      lang: 'tr',
      link: '/tr/',
      themeConfig: {
        nav: [
          { text: 'Rehber', link: '/tr/getting-started' },
          { text: 'API Referansı', link: '/tr/resources/orders' },
          { text: 'Webhook', link: '/tr/webhooks' },
          { text: 'Örnekler', link: 'https://github.com/AdisGroup/shopier-go/tree/main/examples' }
        ],
        sidebar: [
          {
            text: 'Başlangıç',
            items: [
              { text: 'Genel Bakış ve Kurulum', link: '/tr/getting-started' },
              { text: 'Kimlik Doğrulama (PAT & OAuth)', link: '/tr/authentication' },
              { text: 'Sayfalama ve İteratörler', link: '/tr/pagination-and-iterators' },
              { text: 'Hata Yönetimi ve Retry', link: '/tr/error-handling' },
              { text: 'Test ve Mocking', link: '/tr/testing' }
            ]
          },
          {
            text: 'Webhook Entegrasyonu',
            items: [
              { text: 'İmza Doğrulama ve Olaylar', link: '/tr/webhooks' }
            ]
          },
          {
            text: 'API Kaynakları',
            items: [
              { text: 'Siparişler (Orders)', link: '/tr/resources/orders' },
              { text: 'Ürünler (Products)', link: '/tr/resources/products' },
              { text: 'Kategoriler (Categories)', link: '/tr/resources/categories' },
              { text: 'Varyasyonlar (Variations)', link: '/tr/resources/variations' },
              { text: 'Seçenekler (Selections)', link: '/tr/resources/selections' },
              { text: 'İndirimler (Discounts)', link: '/tr/resources/discounts' },
              { text: 'Kargo İşlemleri (Shippings)', link: '/tr/resources/shippings' },
              { text: 'Bakiye (Balance)', link: '/tr/resources/balance' },
              { text: 'Hakedişler (Payouts)', link: '/tr/resources/payouts' },
              { text: 'İadeler (Refunds)', link: '/tr/resources/refunds' },
              { text: 'Mağaza Ayarları (Shop)', link: '/tr/resources/shop' },
              { text: 'Webhook Abonelikleri', link: '/tr/resources/webhooks' }
            ]
          }
        ]
      }
    }
  }
})
