package testdata

const BitAssetData = `
  {
    "id": "2.4.1",
    "feeds": [
      [
        "1.2.6",
        [
          "2018-09-14T16:51:25",
          {
            "settlement_price": {
              "base": {
                "amount": 19424,
                "asset_id": "1.3.2"
              },
              "quote": {
                "amount": 100000000,
                "asset_id": "1.3.0"
              }
            },
            "maintenance_collateral_ratio": 1750,
            "maximum_short_squeeze_ratio": 1500,
            "core_exchange_rate": {
              "base": {
                "amount": 20395,
                "asset_id": "1.3.2"
              },
              "quote": {
                "amount": 100000000,
                "asset_id": "1.3.0"
              }
            }
          }
        ]
      ],
      [
        "1.2.14",
        [
          "2018-09-14T16:52:20",
          {
            "settlement_price": {
              "base": {
                "amount": 19427,
                "asset_id": "1.3.2"
              },
              "quote": {
                "amount": 100000000,
                "asset_id": "1.3.0"
              }
            },
            "maintenance_collateral_ratio": 1750,
            "maximum_short_squeeze_ratio": 1500,
            "core_exchange_rate": {
              "base": {
                "amount": 20399,
                "asset_id": "1.3.2"
              },
              "quote": {
                "amount": 100000000,
                "asset_id": "1.3.0"
              }
            }
          }
        ]
      ],
      [
        "1.2.16",
        [
          "2018-09-14T16:51:25",
          {
            "settlement_price": {
              "base": {
                "amount": 19424,
                "asset_id": "1.3.2"
              },
              "quote": {
                "amount": 100000000,
                "asset_id": "1.3.0"
              }
            },
            "maintenance_collateral_ratio": 1750,
            "maximum_short_squeeze_ratio": 1500,
            "core_exchange_rate": {
              "base": {
                "amount": 20396,
                "asset_id": "1.3.2"
              },
              "quote": {
                "amount": 100000000,
                "asset_id": "1.3.0"
              }
            }
          }
        ]
      ]
    ],
    "current_feed": {
      "settlement_price": {
        "base": {
          "amount": 19426,
          "asset_id": "1.3.2"
        },
        "quote": {
          "amount": 100000000,
          "asset_id": "1.3.0"
        }
      },
      "maintenance_collateral_ratio": 1750,
      "maximum_short_squeeze_ratio": 1500,
      "core_exchange_rate": {
        "base": {
          "amount": 20397,
          "asset_id": "1.3.2"
        },
        "quote": {
          "amount": 100000000,
          "asset_id": "1.3.0"
        }
      }
    },
    "current_feed_publication_time": "2018-09-14T16:51:25",
    "options": {
      "feed_lifetime_sec": 86400,
      "minimum_feeds": 7,
      "force_settlement_delay_sec": 86400,
      "force_settlement_offset_percent": 0,
      "maximum_force_settlement_volume": 2000,
      "short_backing_asset": "1.3.0",
      "extensions": []
    },
    "force_settled_volume": 0,
    "is_prediction_market": false,
    "settlement_price": {
      "base": {
        "amount": 0,
        "asset_id": "1.3.0"
      },
      "quote": {
        "amount": 0,
        "asset_id": "1.3.0"
      }
    },
    "settlement_fund": 0
  }
`
