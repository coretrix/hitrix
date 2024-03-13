"use strict";(self.webpackChunk=self.webpackChunk||[]).push([[9313],{1629:(e,n,s)=>{s.r(n),s.d(n,{data:()=>r});const r={key:"v-1c2050b9",path:"/guide/features/consumer_runners.html",title:"Consumer runners",lang:"en-US",frontmatter:{},excerpt:"",headers:[{level:3,title:"ConsumerRunner - non scalable",slug:"consumerrunner-non-scalable",children:[]},{level:3,title:"ScalableConsumerRunner - scalable",slug:"scalableconsumerrunner-scalable",children:[]}],filePathRelative:"guide/features/consumer_runners.md",git:{updatedTime:1651228845e3,contributors:[{name:"Iliyan",email:"iliyan.motovski@coretrix.com",commits:1}]}}},8660:(e,n,s)=>{s.r(n),s.d(n,{default:()=>l});const r=(0,s(6252).uE)('<h1 id="consumer-runners" tabindex="-1"><a class="header-anchor" href="#consumer-runners" aria-hidden="true">#</a> Consumer runners</h1><p>Consumer runners enable you to quickly spin up BeeORM queue consumers easily. There are 2 types of consumer.</p><ul><li>scalable</li><li>non scalable</li></ul><h3 id="consumerrunner-non-scalable" tabindex="-1"><a class="header-anchor" href="#consumerrunner-non-scalable" aria-hidden="true">#</a> ConsumerRunner - non scalable</h3><p>Use <code>queue.NewConsumerRunner(ctx)</code> to make consumers which are not required to be able to scale.</p><p>This consumer works with following 4 interfaces:</p><ul><li>ConsumerOne (consumes items one by one)</li><li>ConsumerMany (consumes items in batches)</li><li>ConsumerOneByModulo (consumes items one by one using modulo)</li><li>ConsumerManyByModulo (consumes items in batches using modulo)</li></ul><h3 id="scalableconsumerrunner-scalable" tabindex="-1"><a class="header-anchor" href="#scalableconsumerrunner-scalable" aria-hidden="true">#</a> ScalableConsumerRunner - scalable</h3><p>Use <code>queue.NewScalableConsumerRunner(ctx, persistent redis pool)</code> to make consumers which are required to be able to scale.</p><p>This consumer works with following 2 interfaces:</p><ul><li>ConsumerOne (consumes items one by one)</li><li>ConsumerMany (consumes items in batches)</li></ul>',11),u={},l=(0,s(3744).Z)(u,[["render",function(e,n){return r}]])},3744:(e,n)=>{n.Z=(e,n)=>{for(const[s,r]of n)e[s]=r;return e}}}]);