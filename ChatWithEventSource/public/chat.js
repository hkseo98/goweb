$(function () {
  // 웹소켓과 다르게, server-sent 이벤트는 단방향입니다. 데이터 메시지가 서버에서 클라이언트로 (유저의 웹 브라우저 같은) 한 방향으로 전달되는 것입니다. 이 특징은 클라이언트에서 서버로 메시지 형태로 데이터를 보낼 필요가 없을 때, server-sent 이벤트를 훌륭한 선택으로 만든다. 예를 들어,  EventSource 는 소셜 미디어 상태 업데이트, 뉴스피드나 IndexedDB나 web storage같은 클라이언트-사이드 저장 매커니즘으로 데이터를 전달하는 데 유용한 접근법입니다.
  if (!window.EventSource) {
    alert("No EventSource!");
    return;
  }

  var $chatlog = $("#chat-log"); // chat-log id를 가진 태그를 가져옴
  var $chatmsg = $("#chat-msg");

  var isBlank = function (string) {
    return string == null || string.trim() === ""; // ==면 값만 비교, ===면 타입도 비교
  };

  var username;
  while (isBlank(username)) {
    username = prompt("What's your name?");
    if (!isBlank(username)) {
      $("#user-name").html("<b>" + username + "</b>");
    }
  }
  $("#input-form").on("submit", function (e) {
    $.post("/messages", {
      msg: $chatmsg.val(),
      name: username,
    }); // jquery의 post함수를 쓰겠다. - 해당 url로 post 전송
    $chatmsg.val(""); // 전송 후에는 메세지칸을 비우고
    $chatmsg.focus(); // 다시 포커싱
    return false; // false를 반환해야 다른 페이지로 넘어가지 않음
  });

  var addMessage = function (data) {
    var text = "";
    if (!isBlank(data.Name)) {
      text = "<strong>" + data.Name + ":</strong> ";
    }
    text += data.Msg;
    $chatlog.prepend("<div><span>" + text + "</span></div>"); // 앞에서부터 추가
    console.log(data);
  };

  var es = new EventSource("/stream");
  es.onopen = function (e) {
    $.post("/users", {
      name: username,
    });
  };

  es.onmessage = function (e) {
    // es를 통해서 메세지가 올 때
    var msg = JSON.parse(e.data);
    addMessage(msg);
  };

  window.onbeforeunload = function () {
    $.ajax({
      url: "/users?name=" + username,
      type: "DELETE",
    });
    es.close();
  };
});
