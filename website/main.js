(function () {
  "use strict";

  var yearEl = document.getElementById("year");
  if (yearEl) {
    yearEl.textContent = String(new Date().getFullYear());
  }

  var navToggle = document.getElementById("navToggle");
  var siteNav = document.getElementById("siteNav");

  if (navToggle && siteNav) {
    navToggle.addEventListener("click", function () {
      var open = siteNav.classList.toggle("is-open");
      navToggle.setAttribute("aria-expanded", open ? "true" : "false");
    });

    siteNav.querySelectorAll('a[href^="#"]').forEach(function (link) {
      link.addEventListener("click", function () {
        siteNav.classList.remove("is-open");
        navToggle.setAttribute("aria-expanded", "false");
      });
    });
  }

  var form = document.getElementById("contactForm");
  if (!form) return;

  var nameInput = document.getElementById("name");
  var emailInput = document.getElementById("email");
  var messageInput = document.getElementById("message");
  var nameError = document.getElementById("nameError");
  var emailError = document.getElementById("emailError");
  var messageError = document.getElementById("messageError");
  var formSuccess = document.getElementById("formSuccess");

  function clearErrors() {
    if (nameError) nameError.textContent = "";
    if (emailError) emailError.textContent = "";
    if (messageError) messageError.textContent = "";
  }

  function validateEmail(value) {
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value);
  }

  form.addEventListener("submit", function (event) {
    event.preventDefault();
    clearErrors();
    if (formSuccess) formSuccess.hidden = true;

    var valid = true;

    if (!nameInput || !nameInput.value.trim()) {
      if (nameError) nameError.textContent = "请填写姓名";
      valid = false;
    }

    if (!emailInput || !emailInput.value.trim()) {
      if (emailError) emailError.textContent = "请填写工作邮箱";
      valid = false;
    } else if (emailInput && !validateEmail(emailInput.value.trim())) {
      if (emailError) emailError.textContent = "请输入有效的邮箱地址";
      valid = false;
    }

    if (!messageInput || !messageInput.value.trim()) {
      if (messageError) messageError.textContent = "请填写需求说明";
      valid = false;
    }

    if (!valid) return;

    if (formSuccess) {
      formSuccess.hidden = false;
    }
    form.reset();
  });
})();
