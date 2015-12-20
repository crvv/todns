'use strict';

const host = window.location.origin;
const api = host + '/api/';
const apis = {
    upstream: {method: 'PUT', path: 'upstream'},
    ttl: {method: 'PUT', path: 'ttl'},
    add: {method: 'POST', path: 'record'},
    del: {method: 'DELETE', path: 'record'}
};

function addEventListenerList(nodes, event, listener) {
    for (var i = 0; i < nodes.length; i++) {
        nodes[i].addEventListener(event, listener)
    }
}

function request(method, path, body) {
    const req = new XMLHttpRequest();
    req.open(method, api + path);
    req.onload = function onload() {
        if (this.status != 200) {
            return
        }
        window.document.getElementById('config').innerHTML
            = JSON.stringify(JSON.parse(this.responseText), null, '    ');
    };
    req.send(body)
}

addEventListenerList(window.document.getElementsByClassName('button'), 'click', function click() {
    const operation = this.parentNode.id;
    let inputs = this.parentNode.getElementsByClassName("in");
    const body = {};
    for (let i = 0; i < inputs.length; i++) {
        if (inputs[i].classList[1] == 'ttl') {
            body[inputs[i].classList[1]] = +inputs[i].value;
        } else {
            body[inputs[i].classList[1]] = inputs[i].value;
        }
    }
    console.log(body);
    request(apis[operation].method, apis[operation].path, JSON.stringify(body));
});

request('GET', '', '');
