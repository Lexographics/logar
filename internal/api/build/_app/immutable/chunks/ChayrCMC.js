import{am as f,an as l,ao as y,ac as b,M as S,q as m,N as x,ap as B}from"./DOREUks6.js";let w=!1;function M(){w||(w=!0,document.addEventListener("reset",r=>{Promise.resolve().then(()=>{var t;if(!r.defaultPrevented)for(const e of r.target.elements)(t=e.__on_r)==null||t.call(e)})},{capture:!0}))}function E(r){var t=y,e=b;f(null),l(null);try{return r()}finally{f(t),l(e)}}function j(r,t,e,n=e){r.addEventListener(t,()=>E(e));const i=r.__on_r;i?r.__on_r=()=>{i(),n(!0)}:r.__on_r=()=>n(!0),M()}const N=new Set,O=new Set;function T(r,t,e,n={}){function i(a){if(n.capture||W.call(t,a),!a.cancelBubble)return E(()=>e==null?void 0:e.call(this,a))}return r.startsWith("pointer")||r.startsWith("touch")||r==="wheel"?m(()=>{t.addEventListener(r,i,n)}):t.addEventListener(r,i,n),i}function z(r,t,e,n,i){var a={capture:n,passive:i},o=T(r,t,e,a);(t===document.body||t===window||t===document)&&S(()=>{t.removeEventListener(r,o,a)})}function A(r){for(var t=0;t<r.length;t++)N.add(r[t]);for(var e of O)e(r)}function W(r){var g;var t=this,e=t.ownerDocument,n=r.type,i=((g=r.composedPath)==null?void 0:g.call(r))||[],a=i[0]||r.target,o=0,v=r.__root;if(v){var _=i.indexOf(v);if(_!==-1&&(t===document||t===window)){r.__root=t;return}var p=i.indexOf(t);if(p===-1)return;_<=p&&(o=_)}if(a=i[o]||r.target,a!==t){x(r,"currentTarget",{configurable:!0,get(){return a||e}});var L=y,k=b;f(null),l(null);try{for(var s,h=[];a!==null;){var d=a.assignedSlot||a.parentNode||a.host||null;try{var u=a["__"+n];if(u!=null&&(!a.disabled||r.target===a))if(B(u)){var[q,...P]=u;q.apply(a,[r,...P])}else u.call(a,r)}catch(c){s?h.push(c):s=c}if(r.cancelBubble||d===t||d===null)break;a=d}if(s){for(let c of h)queueMicrotask(()=>{throw c});throw s}}finally{r.__root=t,delete r.currentTarget,f(L),l(k)}}}export{N as a,M as b,A as d,z as e,W as h,j as l,O as r,E as w};
