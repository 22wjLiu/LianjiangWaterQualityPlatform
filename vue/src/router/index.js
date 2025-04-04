import Vue from "vue";
import VueRouter from "vue-router";

import Welcome from "@/views/Welcome";
import Register from "@/views/Register";
import Overview from "@/views/Overview";
import Graph from "@/views/Graph";
import File from "@/views/File";
import Login from "@/views/Login";
import BackStage from "@/views/BackStage";
import user from "@/components/backStage/user";
import map from "@/components/backStage/map";
import file from "@/components/backStage/file";
import data from "@/components/backStage/data";
import log from "@/components/backStage/log";
import fileLog from "@/components/backStage/logTables/fileLog";
import dataLog from "@/components/backStage/logTables/dataLog";
import mapLog from "@/components/backStage/logTables/mapLog";
// import { component } from 'vue/types/umd.js'

Vue.use(VueRouter);

const routes = [
  { path: "/", redirect: "/welcome" },
  { path: "/welcome", component: Welcome, meta: { keepAlive: false } },
  { path: "/login", component: Login, meta: { keepAlive: false } },
  { path: "/register", component: Register, meta: { keepAlive: false } },
  { path: "/overview", component: Overview, meta: { keepAlive: true } },
  { path: "/graph", component: Graph, meta: { keepAlive: true } },
  { path: "/file", component: File, meta: { keepAlive: false } },
  { path: "/backStage", redirect: "/backStage/user" },
  {
    path: "/backStage",
    component: BackStage,
    children: [
      { path: "user", component: user },
      { path: "map", component: map },
      { path: "file", component: file },
      { path: "data", component: data },
      {
        path: "log",
        redirect: "/backStage/log/dataLog",
        component: log,
        children: [
          { path: "dataLog", component: dataLog },
          { path: "fileLog", component: fileLog },
          { path: "mapLog", component: mapLog },
        ],
      },
    ],
  },
];

const router = new VueRouter({
  routes,
});

router.beforeEach((to, from, next) => {
  const path = to.path;

  if (
    path === "/graph" ||
    path === "/map" ||
    path === "/file" ||
    path === "/overview"
  ) {
    if (localStorage.getItem("token")) {
      next();
    } else {
      Vue.prototype.$message.warning("请先登录");
      next("/login");
    }
  }
  else if (path === '/backStage/user') {
    if (Number(localStorage.getItem("level")) < 4) {
      Vue.prototype.$message.warning("对不起，您没有权限");
      next(false);
    } else {
      next();
    }
  } 
  else {
    next();
  }

});

export default router;
