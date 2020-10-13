Feature: 更新积分商品

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

	@ginger-promotion @point
	Scenario: 商户能更新禁用的积分商品
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
			"point_price": 9,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后"
		}
		"""
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品1",
			"point_price": 9,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后",
			"is_enabled": true
		}]
		"""

		# jobs更新积分商品
		When jobs禁用积分商品'商品1'
		When jobs更新积分商品'商品1'
		"""
		{
			"point_price": 19,
			"money_price": 3,
			"buy_limit": 2,
			"start_time": "明天",
			"end_time": "7天后"
		}
		"""
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品1",
			"point_price": 19,
			"money_price": 3,
			"buy_limit": 2,
			"start_time": "明天",
			"end_time": "7天后",
			"is_enabled": false
		}]
		"""

	@ginger-promotion @point
	Scenario: 商户不能更新启用的积分商品
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
			"point_price": 9,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后"
		}
		"""
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品1",
			"point_price": 9,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后",
			"is_enabled": true
		}]
		"""

		# jobs更新积分商品
		When jobs更新积分商品'商品1'
		"""
		{
			"point_price": 19,
			"money_price": 3,
			"buy_limit": 2,
			"start_time": "明天",
			"end_time": "7天后",
			"error": "product:update_enabled_point_product"
		}
		"""
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品1",
			"point_price": 9,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后",
			"is_enabled": true
		}]
		"""


