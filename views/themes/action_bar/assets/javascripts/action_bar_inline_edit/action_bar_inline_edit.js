(function (factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as anonymous module.
    define(['jquery'], factory);
  } else if (typeof exports === 'object') {
    // Node / CommonJS
    factory(require('jquery'));
  } else {
    // Browser globals.
    factory(jQuery);
  }
})(function ($) {

  'use strict';

  var $body = $("body");
  var NAMESPACE = 'qor.actionbar.inlineEdit';
  var EVENT_ENABLE = 'enable.' + NAMESPACE;
  var EVENT_DISABLE = 'disable.' + NAMESPACE;
  var EVENT_CLICK = 'click.' + NAMESPACE;
  var EDIT_ACTIONBAR_BUTTON = '.qor-actionbar-button';
  var ID_ACTIONBAR = 'qor-actionbar-iframe';
  var INLINE_EDIT_URL;

  function QorActionBarInlineEdit(element, options) {
    this.$element = $(element);
    this.options = $.extend({}, QorActionBarInlineEdit.DEFAULTS, $.isPlainObject(options) && options);
    this.init();
  }

  QorActionBarInlineEdit.prototype = {
    constructor: QorActionBarInlineEdit,

    init: function () {
      this.bind();
      this.initStatus();
    },

    bind: function () {
      this.$element.on(EVENT_CLICK, EDIT_ACTIONBAR_BUTTON, this.click);
      $(document).on('keyup', this.keyup);
    },

    initStatus : function () {
      var iframe = document.createElement("iframe");

      iframe.src = INLINE_EDIT_URL;
      iframe.id = ID_ACTIONBAR;

      // show edit button after iframe totally loaded.
      if (iframe.attachEvent){
        iframe.attachEvent("onload", function(){
          $('.qor-actionbar-button').show();
        });
      } else {
        iframe.onload = function(){
          $('.qor-actionbar-button').show();
        };
      }

      document.body.appendChild(iframe);
    },

    keyup: function (e) {
      var iframe = document.getElementById('qor-actionbar-iframe');
      if (e.keyCode === 27) {
        iframe && iframe.contentDocument.querySelector('.qor-slideout__close').click();
      }
    },

    click: function () {
      var $this = $(this);
      var iframe = document.getElementById('qor-actionbar-iframe');
      var $iframe = iframe.contentWindow.$;
      var editLink = iframe.contentDocument.querySelector('.js-actionbar-edit-link');

      if (!editLink) {
        return;
      }

      iframe.classList.add('show');

      if ($iframe) {
        $iframe(".js-actionbar-edit-link").attr("data-url", $this.data("url")).click();
      } else {
        editLink.setAttribute("data-url", $this.data("url"));
        editLink.click();
      }

      $body.addClass("open-actionbar-editor");

      return false;
    }
  };

  QorActionBarInlineEdit.plugin = function (options) {
    return this.each(function () {
      var $this = $(this);
      var data = $this.data(NAMESPACE);
      var fn;

      if (!data) {
        if (/destroy/.test(options)) {
          return;
        }
        $this.data(NAMESPACE, (data = new QorActionBarInlineEdit(this, options)));
      }

      if (typeof options === 'string' && $.isFunction(fn = data[options])) {
        fn.apply(data);
      }
    });
  };


  $(function () {
    $body.attr("data-toggle", "qor.actionbars");

    var selector = '[data-toggle="qor.actionbars"]';

    $(".qor-actionbar-button").each(function(i, e) {
      INLINE_EDIT_URL = $(e).data("iframe-url");
    });

    $(document).
      on(EVENT_DISABLE, function (e) {
        QorActionBarInlineEdit.plugin.call($(selector, e.target), 'destroy');
      }).
      on(EVENT_ENABLE, function (e) {
        QorActionBarInlineEdit.plugin.call($(selector, e.target));
      }).
      triggerHandler(EVENT_ENABLE);
  });

  return QorActionBarInlineEdit;
});
