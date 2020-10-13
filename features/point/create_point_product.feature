Feature: 创建积分商品

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
	Scenario: 商户创建积分商品
		# jobs创建商品
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

		# jobs初始验证
		Then jobs能获得积分商品列表
		"""
		[]
		"""

		# jobs创建积分商品
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
		When jobs添加积分商品'商品3'
		"""
		{
			"point_price": 19,
			"money_price": 1,
			"buy_limit": 2
		}
		"""
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品3",
			"point_price": 19,
			"money_price": 1,
			"buy_limit": 2,
			"is_enabled": true
		}, {
			"name": "商品1",
			"point_price": 9,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后",
			"is_enabled": true
		}]
		"""

		# bill验证
		Given bill登录系统
		Then bill能获得积分商品列表
		"""
		[]
		"""

		# lucy验证jobs的积分商品
		Given lucy访问'jobs'的商城
		Then lucy能在商城中看到积分商品列表
		"""
		[{
			"name": "商品3",
			"point_price": 19,
			"money_price": 1,
			"buy_limit": 2,
			"is_enabled": true
		}, {
			"name": "商品1",
			"point_price": 9,
			"money_price": 0,
			"buy_limit": 1,
			"start_time": "今天",
			"end_time": "5天后",
			"is_enabled": true
		}]
		"""

		# lucy验证bill的积分商品
		Given lucy访问'bill'的商城
		Then lucy能在商城中看到积分商品列表
		"""
		[]
		"""


