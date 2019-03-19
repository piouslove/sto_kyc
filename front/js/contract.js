var KYCAbi = [
	{
		"constant": false,
		"inputs": [
			{
				"name": "_newOwner",
				"type": "address"
			}
		],
		"name": "setOwner",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_who",
				"type": "address"
			}
		],
		"name": "certify",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_who",
				"type": "address"
			}
		],
		"name": "getCertifier",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "string"
			}
		],
		"name": "getAddress",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "delegates",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_old",
				"type": "address"
			}
		],
		"name": "removeDelegate",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_who",
				"type": "address"
			}
		],
		"name": "revoke",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_whos",
				"type": "address[]"
			}
		],
		"name": "certifys",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "string"
			}
		],
		"name": "getUint",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_who",
				"type": "address"
			}
		],
		"name": "certified",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_new",
				"type": "address"
			}
		],
		"name": "addDelegate",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "string"
			}
		],
		"name": "get",
		"outputs": [
			{
				"name": "",
				"type": "bytes32"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "who",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "by",
				"type": "address"
			}
		],
		"name": "Confirmed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "who",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "by",
				"type": "address"
			}
		],
		"name": "Revoked",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "who",
				"type": "address"
			}
		],
		"name": "Confirmed",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "who",
				"type": "address"
			}
		],
		"name": "Revoked",
		"type": "event"
	}
];

var KYCAddress = "0xc740d1959803a61e5c89273a4854da2efec28375";

var etherscanProvider = new ethers.providers.EtherscanProvider('ropsten');

var contractInstance = new ethers.Contract(KYCAddress, KYCAbi, etherscanProvider);

var server_url = "http://localhost:2223";

var managerWallet;

var managerMsg;
var managerSig;

var contract;

var checkedData = {
	id: 0, 
	address: "",
	name: "",
	email: "",
	selector: "",
	passport: ""
};

var cardHtml1 = '<div class="card text-center bg-secondary"><div class="card-header">待认证信息</div><div class="card-body"><h6 class="card-title">姓名：'
var cardHtml2 = '</h6><h6 class="card-title">Email：'
var cardHtml3 = '</h6><h6 class="card-title">地址：'
var cardHtml4 = '</h6><picture><source type="image"><img src="'
var cardHtml5 = '" class="img-fluid img-thumbnail" ></picture><br><br><a class="btn btn-primary" type="button" onclick="certify()">通过</a><a>&nbsp;&nbsp;&nbsp;&nbsp;</a><a class="btn btn-primary" type="button" onclick="reject()">拒绝</a></div><div class="card-footer text-muted">data from idhub</div></div>'
