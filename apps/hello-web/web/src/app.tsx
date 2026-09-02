import React from 'https://esm.sh/react@19.1.1';
import {createRoot} from 'https://esm.sh/react-dom@19.1.1/client';

const h=React.createElement;
let root=globalThis.__playgroundRoot;

function App(){
  const [route,setRoute]=React.useState('');
  const [candidate,setCandidate]=React.useState('');
  const [diagnostic,setDiagnostic]=React.useState(null);
  React.useEffect(()=>{fetch('/__dev/current').then(r=>r.json()).then(v=>{setRoute(v.route||'');setCandidate(v.route||'')});fetch('/api/diagnostic').then(r=>r.json()).then(setDiagnostic).catch(e=>setDiagnostic({error:String(e)}))},[]);
  return h(React.Fragment,null,
    h('header',null,h('span',{className:'eyebrow'},'Bazel × Kubernetes'),h('h1',null,'Dev CLI Playground'),h('p',null,'The selected route independently overlays each service and falls back to base.')),
    h('section',{className:'menu'},h('label',null,'Branch route key',h('input',{value:candidate,onChange:e=>setCandidate(e.target.value),placeholder:'owner-feature-hash'})),h('div',null,h('button',{onClick:()=>location.assign('/__dev/select?branch='+encodeURIComponent(candidate))},'Select branch'),h('button',{className:'secondary',onClick:()=>location.assign('/__dev/select')},'Clear')),h('p',null,'Selected: ',h('code',null,route||'base'))),
    h('section',null,h('h2',null,'Request diagnostics'),h('pre',null,JSON.stringify(diagnostic,null,2)))
  );
}
export async function render(){if(!root){root=createRoot(document.getElementById('root'));globalThis.__playgroundRoot=root}root.render(h(App));}
