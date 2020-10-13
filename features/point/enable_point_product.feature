Feature: 启用禁用积分商品

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
	Scenario: 商户启用禁用积分商品
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
			"point_price": 9
		}
		"""
		When jobs添加积分商品'商品2'
		"""
		{
			"point_price": 19
		}
		"""
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品2",
			"point_price": 19,
			"is_enabled": true
		}, {
			"name": "商品1",
			"point_price": 9,
			"is_enabled": true
		}]
		"""

		# lucy验证jobs的积分商品
		Given lucy访问'jobs'的商城
		Then lucy能在商城中看到积分商品列表
		"""
		[{
			"name": "商品2",
			"point_price": 19
		}, {
			"name": "商品1",
			"point_price": 9
		}]
		"""

		# jobs禁用积分商品
		Given jobs登录系统
		When jobs禁用积分商品'商品1'
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品2",
			"point_price": 19,
			"is_enabled": true
		}, {
			"name": "商品1",
			"point_price": 9,
			"is_enabled": false
		}]
		"""

		# lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能在商城中看到积分商品列表
		"""
		[{
			"name": "商品2",
			"point_price": 19
		}]
		"""

		# jobs启用积分商品
		Given jobs登录系统
		When jobs启用积分商品'商品1'
		Then jobs能获得积分商品列表
		"""
		[{
			"name": "商品2",
			"point_price": 19,
			"is_enabled": true
		}, {
			"name": "商品1",
			"point_price": 9,
			"is_enabled": true
		}]
		"""

		# lucy验证
		Given lucy访问'jobs'的商城
		Then lucy能在商城中看到积分商品列表
		"""
		[{
			"name": "商品2",
			"point_price": 19
		}, {
			"name": "商品1",
			"point_price": 9
		}]
		"""


