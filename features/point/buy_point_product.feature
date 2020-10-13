Feature: 购买积分商品

	Background:
		Given 系统配置虚拟资产
		"""
		[{
			"code": "point",
			"display_name": "积分",
			"exchange_rate": 1,
			"enable_fraction": false,
			"is_payable": true,
			"is_debtable": false
		}]
		"""

		Given ginger登录系统
		When ginger创建公司
		"""
		[{
			"name": "MIX",
			"username": "jobs"
		}, {
			"name": "BabyFace",
			"username": "bill"
		}]
		"""

	@ginger-promotion @point @wip
	Scenario: 用户购买积分商品
		# jobs创建积分商品
		Given jobs登录系统
		When jobs添加商品
		"""
		[{
			"name": "商品1",
			"price": 4.00
		}, {
			"name": "商品2",
			"price": 5.00
		}, {
			"name": "商品3",
			"price": 6.00
		}]
		"""
		When jobs添加积分商品'商品1'
		"""
		{
			"point_price": 8,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后"
		}
		"""
		When jobs添加积分商品'商品3'
		"""
		{
			"point_price": 19,
			"money_price": 1,
			"buy_limit": 2
		}
		"""

		# lucy购买jobs的积分商品
		Given lucy访问'jobs'的商城
		When lucy充值'9'个'point'
		When lucy购买'jobs'的商品
		"""
		{
			"products": [{
				"name": "商品1",
				"count": 1
			}],
			"imoneys": [{
				"code": "point",
				"count": 7
			}]
		}
		"""
		Then lucy成功创建订单
		"""
		{
			"status": "待支付",
			"final_money": 1.0,
			"delivery_items": [{
				"final_money": 1.0,
				"product_price": 8.0,
				"imoneys": [{
					"code": "point",
					"count": 7,
					"deduction_money": 7.0
				}]
			}],
			"imoneys": [{
				"code": "point",
				"count": 7,
				"deduction_money": 7.0
			}]
		}
		"""
		Then lucy能获得虚拟资产'point'
		"""
		{
			"balance": 2
		}
		"""



