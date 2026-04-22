(function () {
  var src = new EventSource("/events");
  src.onmessage = function (e) {
    if (e.data === "reload") {
      window.location.reload();
    }
  };
  src.onerror = function () {
    setTimeout(function () {
      if (src.readyState === EventSource.CLOSED) {
        src = new EventSource("/events");
      }
    }, 1000);
  };
})();
