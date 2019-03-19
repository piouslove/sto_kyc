function getKYCItems() {
	$.ajax({
		url : server_url + "/KYCItems",
		async : true, 
		type : "GET", 
		success : function(res) {
			var html = "";
			for (var i = 0; i <= res.length - 1; i++) {
				html += '<option>' + res[i] + '</option>'
			}
			document.getElementById("formSelect").innerHTML = html;
		}, 
		dataType : "json"
	});
}

function applyServer() {
	var formData = new FormData();
	formData.append("address", document.getElementById("formAddress").value);
	// formData.append("name", document.getElementById("formName").value);
	formData.append("email", document.getElementById("formEmail").value);
	formData.append("name", document.getElementById("formName").value);
	formData.append("selector", document.getElementById("formSelect").value);
	formData.append("passport", document.getElementById("formFile").files[0]);
	$.ajax({
		url : server_url + "/apply",
		async : false, 
		processData: false,
		contentType: false,
		type : "POST", 
		data : formData, 
		success : function(res) {
			alert(res.success);
		}, 
		error : function(res) {
			alert(res.responseText);
		}, 
		dataType : "json"
	});
}

function queryServer() {
	var formData = new FormData();
	formData.append("address", document.getElementById("formAddress").value);
	formData.append("selector", document.getElementById("formSelect").value);
	$.ajax({
		url : server_url + "/query",
		async : false, 
		processData: false, 
		contentType: false, 
		type : "POST", 
		data : formData, 
		success : function(res) {
			alert(res.success);
		}, 
		error : function(res) {
			alert(res.responseText);
		},
		dataType : "json"
	});
}

function applyEthereum() {
	var address = document.getElementById("formAddress").value;
	console.log("开始链上")
	return contractInstance.certified(address);
}

function queryEthereum() {
	var address = document.getElementById("formAddress").value;
	return contractInstance.certified(address);
}

function apply() {
	var address = document.getElementById("formAddress").value;
	console.log("开始链上")
	contractInstance.certified(address).then(function(r) {
		console.log(r);
		if (r == true) {
			alert("此地址已被认证过！")
		} else {
			applyServer();
		}
	}).catch(function(r){
		alert("链上查询失败，请重试！");
	});
}

function query() {
	queryEthereum().then(function(r) {
		if (r == true) {
			queryServer();
			// alert("此地址已被认证过！")
		} else {
			alert("还未审核通过，请稍后再试！")
		}
	}).catch(function(r){
		alert("链上查询失败，请重试！")
	});
}

function initWallet() {
	var key = document.getElementById("formPassword").value;
	var file = document.getElementById("formKeyStore").files[0];
	var reader = new FileReader();
	reader.readAsText(file, "UTF-8");
	var info = "";
	reader.onload = function() { 
   		info = reader.result;
   		_initWallet(info, key);
	};
}

function _initWallet(info, key) {

	// let json = JSON.stringify(info);
	// let wallet;
	ethers.Wallet.fromEncryptedJson(info, key).then(function(r) {
		wallet = r;
    	console.log("Address: " + wallet.address);
    	var walletWithProvider = wallet.connect(etherscanProvider);
		contract = contractInstance.connect(walletWithProvider);
		var msg = Number(new Date()).toString();
		wallet.signMessage(msg).then(function(r){
			managerMsg = msg;
			managerSig = r;
			getDataToCheck();
		}).catch(function(r){
			alert("加密身份验证出错，请检查您的输入信息！");
		});
    	}).catch(function(r) {
    		alert("KeyStore文件或密码输入错误，请检查您的输入信息！");
	});
}

function getDataToCheck() {
	checkedData.selector = document.getElementById("formSelect").value;
	_getDataToCheck();
}

function _getDataToCheck() {
	console.log(checkedData);
	var formData = new FormData();
	formData.append("selector", checkedData.selector);
	formData.append("userId", checkedData.id);
	formData.append("msg", managerMsg);
	formData.append("sig", managerSig);
	$.ajax({
		url : server_url + "/getDataToCheck",
		async : false, 
		processData: false, 
		contentType: false, 
		type : "POST", 
		data : formData, 
		success : function(res) {
			checkedData = {
				id: res.id, 
				address: res.address,
				name: res.name,
				selector: res.selector,
				passport: res.passport
			};
			var html = cardHtml1 + res.name + cardHtml2 + res.email + cardHtml3 + res.address;
			html += cardHtml4 + server_url + "/passportImages/" + res.passport + cardHtml5;
			document.getElementById("checkedData").innerHTML =  html;
		}, 
		error : function(res) {
			document.getElementById("checkedData").innerHTML =  '<br><br><h3>暂无待认证用户！<h3><br><br>';
			alert(res.responseText);
		},
		dataType : "json"
	});
	console.log(checkedData);
}

function certify() {
	var formData = new FormData();
	formData.append("userId", checkedData.id);
	formData.append("msg", managerMsg);
	formData.append("sig", managerSig);
	contract.certify(checkedData.address).then(function(r){
		console.log(r)
		$.ajax({
			url : server_url + "/certify",
			async : false, 
			processData: false, 
			contentType: false, 
			type : "POST", 
			data : formData, 
			success : function(res) {
				console.log(res.success)
				alert(res);
			}, 
			error : function(res) {
				alert(res.responseText);
			},
			dataType : "json"
		});
		_getDataToCheck();
	}).catch(function(res){
		console.log(res);
		alert("此地址已上链！");
		_getDataToCheck();
	})
	console.log(checkedData);
}

function reject() {
	var formData = new FormData();
	formData.append("userId", checkedData.id);
	formData.append("msg", managerMsg);
	formData.append("sig", managerSig);
	$.ajax({
		url : server_url + "/reject",
		async : false, 
		processData: false, 
		contentType: false, 
		type : "POST", 
		data : formData, 
		success : function(res) {
			console.log(res)
			alert(res.success);
		}, 
		error : function(res) {
			alert(res.responseText);
		},
		dataType : "json"
	});
	_getDataToCheck();
	console.log(checkedData);
}

function init() {
	getKYCItems();
}

init();























