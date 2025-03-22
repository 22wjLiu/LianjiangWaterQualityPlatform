# lianjiang

- ## 用户相关

  - **接口地址：/regist**

    **功能：用于用户注册**

    **方法类型：POST**

    请求参数：Body部分，form-data类型，接收四个字符串分别为Email，Password，Name，Verify。其中Email需要符合邮箱格式，Password最少需要六位，Name最多为20位长度，不输入name会生成一段随机字符串，Verify必须与邮箱验证码相同，注意，用户的邮箱和名称均不能重复。

    可能的返回值：成功则返回200与token，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/verify/:id**

    **功能：用于用户请求验证邮箱**

    **方法类型：GET**

    请求参数：需要在接口地址部分（:id）给出用户邮箱

    返回值：成功则返回200并向相应邮箱发送验证邮件，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/security**

    **功能：用于用户找回密码**

    **方法类型：PUT**

    请求参数：Body部分，form-data类型，接收两个字符串分别为Email，Verify。其中Verify必须与邮箱验证码相同。

    返回值：成功则返回200并向相应邮箱发送重置后的密码，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/updatepass**

    **功能：用于用户修改密码**

    **方法类型：PUT**

    请求参数：Body部分，form-data类型，接收两个字符串分别为first，second。其中first为旧密码，second为新密码。

    返回值：成功则返回200并修改数据库中的密码，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/login**

    **功能：用于用户登录**

    **方法类型：POST**

    请求参数：Body部分，form-data类型，接收两个字符串分别为Email，Password。其中Email需要符合邮箱格式，Password最少需要六位。

    可能的返回值：成功则返回200与token，失败则返回其他值，msg中将会给出失败原因

  - **接口地址：/personal**

    **功能：获取当前用户的个人信息**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：返回用户的一些个人信息，格式为json包含name,email,level，name表示用户名，email表示用户邮箱，level表示用户权限等级

  - **接口地址：/level/:id/:level**

    **功能：设置特定用户的用户等级**

    **方法类型：PUT**

    请求参数：需要当前登录用户拥有5及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token。被修改用户的id（接口地址中的id），需要修改的等级level(接口地址中的level)

    返回值：返回设置成功信息

  - **接口地址：/users**

    **功能：获取用户列表**

    **方法类型：GET**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token。

    返回值：返回设一个users，其为user数组，每个user包含name,email,level，name表示用户名，email表示用户邮箱，level表示用户权限等级

  - **接口地址：/user/:search/:id**

    **功能：获取指定用户的信息**

    **方法类型：GET**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，给出需要搜索的字段search，和字段对应的值id。search的允许字段为"名称"，邮箱"，"ID"。

    返回值：返回一个use包含name,email,level，name表示用户名，email表示用户邮箱，level表示用户权限等级

- ## 映射相关

  允许操作的映射表以及初始值如下：

  ````go
  "列字段映射":
  {
      ("监测断面", "StationName")
  	("监测指标", "Time")
  	("监测时间", "Time")
  	("时间", "Time")
  	("水温", "Temperature")
  	("pH", "PH")
  	("化学需氧量", "Cod")
  	("五日生化需氧量", "FiveDaysNiochemicalOxygenDemand")
  	("硒", "Se")
  	("砷", "As")
  	("汞", "Hg")
  	("氟化物", "Fluoride")
  	("石油类", "Petroleum")
  	("粪大肠菌群", "FecalColiform")
  	("溶解氧", "DO")
  	("电导率", "EC")
  	("浊度", "Turbidity")
  	("高锰酸盐指数", "CODMII")
  	("氨氮", "NH_N")
  	("总磷", "TP")
  	("总氮", "TN")
  	("CODcr", "CODcr")
  	("氰化物", "CN")
  	("挥发酚", "VolatilePenol")
  	("六价铬", "Cr")
  	("铜", "Cu")
  	("锌", "Zn")
  	("铅", "Pb")
  	("镉", "Cd")
  	("阴离子表面活性剂", "LAS")
  	("硫化物", "SOx")
  	("累计流量", "CumulativeDischarge")
  	("水流量", "WaterDischarge")
  	("总累积流量", "TotalCumulativeFlow")
  	("水位", "WaterLevel")
  	("时段累积流量", "PeriodCumulativeFlow")
  	("断面平均流速", "SectionalMeanVelocity")
  	("当前瞬时流速", "CurrentInstantaneousVelocity")
  	("瞬时流量", "InstantaneousDelivery")
  	("断面面积", "SectionalArea")
  }
  
  "行字段一对多映射":
  {
      ("分项类别", "item_category")
  }
  
  "行字段一对一映射":
  {
      ("水质类别", "water_quality_classification")
  	("主要污染物", "key_pollutant")
  }
  
  "时间制映射":
  {
      ("小时制", "hour")
  	("月度制", "month")
  }
  
  "标记映射":
  {
      ("hour", "时间")
  	("month", "监测断面")
  }
  
  "站名映射":
  {
      ("海门湾桥闸", "haimen_bay_bridge_gate")
  	("汕头练江水站", "lian_jiang_water_station")
  	("青洋山桥", "lian_jiang_water_station")
  	("新溪西村", "xinxi_village")
  	("万兴桥", "wanxing_bridge")
  	("流仙学校", "liuxian_school")
  	("仙马闸", "xianma_brake")
  	("华侨学校", "huaqiao_school")
  	("港洲桥", "gangzhou_bridge")
  	("云陇", "yunlong")
  	("北港水", "beixiangshui")
  	("官田水", "guantianshui")
  	("北港河闸", "beixiang_penstock")
  	("峡山大溪", "xiashan_stream")
  	("井仔湾闸", "jingzai_wan_sluice")
  	("东北支流", "northeast_branch")
  	("西埔桥闸", "xipu_bridge_sluice")
  	("五福桥", "wufu_bridge")
  	("成田大寮", "narita_daliao")
  	("新坛港", "xitan_port")
  	("瑶池港", "yaochi_port")
  	("护城河闸", "moat_locks")
  	("和平桥", "peace_bridge")
  }
  
  "数据符号映射":
  {
      
  }
  ````

  - **接口地址：/map/:id**

    **功能：查看指定映射表的所有key**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token，在id处给出表名

    返回值：返回一个keys，为string数组，包含该表的所有key

  - **接口地址：/map/:id/:key**

    **功能：查看指定映射表的指定key的值**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token，在id处给出表名，key处给出指定的key

    返回值：返回一个value，表示key对应的value

  - **接口地址：/map/:id**

    **功能：通过同名键值更新或创建映射**

    **方法类型：PUT**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在id处给出表名，在Params处提供key1（表示用于取值的键值）和key2（表示新建的键值），将会创建与key1同值的key2。

    返回值：返回创建成功信息

  - **接口地址：/map/:id/:key**

    **功能：更新或创建映射键值对**

    **方法类型：PUT**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在id处给出表名， 在key处给出要创建或修改的键值，在Params处提供value(默认为空)

    返回值：返回创建成功信息

  - **接口地址：/map/:id/:key**

    **功能：删除映射键值对**

    **方法类型：DELETE**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在id处给出表名， 在key处给出要删除的键值

    返回值：返回删除成功信息

- ## 文件相关

  - **接口地址：/upload/:system**

    **功能：文件上传**

    **方法类型：POST**

    请求参数：需要当前登录用户拥有3及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在system处给出时间制度，其必须在**时间制映射**中存在，并在**标记映射**中存在标记。在Body部分，form-data格式，接收file（文件类型），先仅支持.xls,.xlsx,.csv文件，大小不超过10M

    返回值：返回上传成功信息

  - **接口地址：/files**

    **功能：获取文件列表**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token，在Params处提供path，表示要查找的文件路径

    返回值：返回files，为file数组，每个file包含name、path、type、lastWriteTime、size，name表示文件名、path表示文件路径、type表示文件类型、lastWriteTime表示文件最后修改日期、size表示文件大小

  - **接口地址：/download**

    **功能：下载指定文件**

    **方法类型：GET**

    请求参数：需要当前登录用户拥有2及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在Params处提供path和file，path表示要下载的文件路径，file表示下载的文件的文件名

    返回值：返回指定文件

  - **接口地址：/file**

    **功能：文件删除**

    **方法类型：DELETE**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在Params处提供path和file，path表示要删除的文件路径

    返回值：返回删除成功信息

- ## 数据相关

  - **接口地址：/data/:name/:system**

    **功能：数据获取**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token，在name处给出站名，必须在**站名映射**中存在，在system处给出时间制，必须在**时间制映射**中存在，在Params处提供fields，start和end，start和end表示初始时间和终止时间，fields表示field数组，每个field表示需要查找的字段，必须在**列字段映射**中存在，当start为空值则不受start限制，end为空值则不受end限制

    返回值：返回resultArr数组，表示在start和end时间之间的表中数据，其每个元素包含time和field，分别表示记录时间和字段的值

  - **接口地址：/data/rowall/:key/:name**

    **功能：获取一对多的行字段数据**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token，在key处给出一对多的主键，必须在**行字段一对多映射**中存在，在name处给出站名，必须在**站名映射**中存在，在Params处提供fields、start和end，fields表示一个数组，其元素必须在**列字段映射**存在，表示初始时间和终止时间，当start为空值则不受start限制，end为空值则不受end限制

    返回值：返回resultArr数组，表示在start和end时间之间的表中数据，其每个元素包含start_time、end_time和field，分别表示记录时间段和字段的值

  - **接口地址：/data/rowone/:key/:name**

    **功能：获取一对一的行字段数据**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token，在key处给出一对一的主键，必须在**行字段一对一映射**中存在，在name处给出站名，必须在**站名映射**中存在，在Params处提供start和end，表示初始时间和终止时间，当start为空值则不受start限制，end为空值则不受end限制

    返回值：返回resultArr数组，表示在start和end时间之间的表中数据，其每个元素包含start_time、end_time和field，分别表示记录时间段和字段的值

  - **接口地址：/data/:time/:start/:end**

    **功能：数据删除**

    **方法类型：DELETE**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在time处给出用于参考的时间字段，**只能为"创建日期"或"记录日期"**，在start处给出起始时间，在end处给出终止时间，当start为空值则不受start限制，end为空值则不受end限制。在Params处提供system和name，system必须在**时间制映射**中存在，name必须在**站名映射**中存在。如果system为空则默认操作所有时间制，name为空则默认操作所有站名。将会删除对应时间制站名表的time在start和end之间的所有数据。

    返回值：返回数据删除成功信息

  - **接口地址：/data/:start/:end**

    **功能：数据恢复**

    **方法类型：PUT**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在start处给出起始时间，在end处给出终止时间，当start为空值则不受start限制，end为空值则不受end限制。在Params处提供system和name，system必须在**时间制映射**中存在，name必须在**站名映射**中存在。如果system为空则默认操作所有时间制，name为空则默认操作所有站名。将会恢复对应时间制站名表的在起始时间和终止时间删除的所有数据。

    返回值：返回数据恢复成功信息
    
  - **接口地址：/forecast**

    **功能：数据预测**

    **方法类型：GET**

    请求参数：Authorization中的Bearer Token中提供注册、登录时给出的token，在Params中给出Temperature、PH、Turbidity、DO四个参数，每个参数均为浮点数。

    返回值：返回一个data数组，包含了三个浮点数

- ## 记录相关

  - **接口地址：/history/file/:start/:end**

    **功能：查看文件上传、删除记录**

    **方法类型：GET**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在start处给出起始时间，在end处给出终止时间，当start为空值则不受start限制，end为空值则不受end限制。在Params处提供id（可选）、option（可选），id表示用户id，option表示操作类型，如不为空则必须为"创建"或"删除"。

    返回值：返回fileHistorys，其为fileHistory数组，每个fileHistory包含了user_id、created_at、file_name、file_path、option，分别表示用户id，历史记录创建日期，文件名，文件路径，用户操作

  - **接口地址：/history/data/:start/:end**

    **功能：查看数据上传、删除记录**

    **方法类型：GET**

    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在start处给出起始时间，在end处给出终止时间，当start为空值则不受start限制，end为空值则不受end限制。在Params处提供id（可选）、option（可选）、station_name（可选）、system（可选），id表示用户id，option表示操作类型，如不为空则必须为"创建"或"删除"，station_name表示站名，如不为空则必须在**站名映射**中存在，system表示时间制，如果为空则必须在**时间制映射**中存在。

    返回值：返回dataHistorys，其为dataHistory数组，每个dataHistory包含了user_id、created_at、start_time、end_time、option、station_name、system、time，user_id表示用户id，created_at表示历史记录创建日期，start_time，end_time和time表示操作的时间区间为参考字段time大于等于start_time小于等于end_time，station_name表示站名，system表示时间制
    
  - **接口地址：/history/map/:start/:end**
  
    **功能：查看用户的映射上传、删除记录**
  
    **方法类型：GET**
  
    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在start处给出起始时间，在end处给出终止时间，当start为空值则不受start限制，end为空值则不受end限制。在Params处提供id（可选）、option（可选）、station_name（可选）、system（可选），id表示表名，option表示操作类型，如不为空则必须为"创建"或"删除"，user_id表示操作者名称。
  
    返回值：返回mapHistorys，其为mapHistory数组，每个mapHistory包含了id、user_id、created_at、option、key、value，id表示表名称，user_id表示用户id，created_at表示历史记录创建日期，key表示主键，value表示值。
  
  - **接口地址：/map/backup/:id/:start/:end**
  
    **功能：查看映射的备份**
  
    **方法类型：GET**
  
    请求参数：需要当前登录用户拥有4及以上安全等级。Authorization中的Bearer Token中提供注册、登录时给出的token，在id处给出表名，在start处给出起始时间，在end处给出终止时间。
  
    返回值：返回backupMaps，其为backupMap数组，每个backupMap包含了create_at,id,key,value,其中create_at表示备份的创建时间，table表示映射表名，key表示主键，value表示值。
    
    