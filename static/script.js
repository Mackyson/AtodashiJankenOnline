function startChat(ws){
	//入室時の発言
	ws.addEventListener("open",function(e){
		console.log("WebSocket connected")
		var data = {}
		data["name"] = "System"
		ws.send(
			JSON.stringify(data)
		)
	});
	//発言を受信したらlistに表示
	ws.addEventListener("message",function(e){
		json = e.data
		msg = JSON.parse(json)
		console.log(msg)
		// var user=msg["name"],payload=msg["payload"],time=msg["time"]
		// var li = document.createElement("li");
		// li.textContent=user+" : "+payload+" ("+time+")"
		// 	document.getElementById("list").appendChild(li);
	});
	//boxに書いた内容を発言
	document.getElementById("sendBtn").addEventListener("click",function(e){
		var box = document.getElementById("box")
		var data = {}
		data["time"] = ""
		ws.send(
			JSON.stringify(data)
		)
		box.value = ""
	});
}
function enter(){
	//displayの値で画面を切り替え
	entrance.style.display="none";
	chat.style.display="block";
	//エントリー用のハンドラに接続
	ws = new WebSocket("ws://localhost:8080/match")
	ws.addEventListener("open",function(e){
	startChat(ws)
	})
}
