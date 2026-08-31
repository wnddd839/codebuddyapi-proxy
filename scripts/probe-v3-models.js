const fs = require("fs");
const https = require("https");

const os = require("os");
const path = require("path");

const defaultPath = path.join(os.homedir(), ".codebuddy", "proxy-accounts.json");
const accountsPath = process.argv[2] || defaultPath;
const store = JSON.parse(fs.readFileSync(accountsPath, "utf8"));
const IDE = "2.117.2";

function req(url, token) {
  return new Promise((resolve, reject) => {
    const u = new URL(url);
    const headers = {
      Accept: "application/json",
      Authorization: "Bearer " + token,
      "X-IDE-Type": "CLI",
      "X-IDE-Name": "CLI",
      "X-IDE-Version": IDE,
      "User-Agent": "CLI/" + IDE + " CodeBuddy/" + IDE,
      "X-Product": "SaaS",
      "X-User-Id": "probe",
      "X-Domain": u.hostname,
      "X-Agent-Intent": "craft",
      "X-Requested-With": "XMLHttpRequest",
    };
    const r = https.request({ hostname: u.hostname, path: u.pathname, method: "GET", headers }, (res) => {
      let d = "";
      res.on("data", (c) => (d += c));
      res.on("end", () => resolve({ status: res.statusCode, body: d }));
    });
    r.on("error", reject);
    r.end();
  });
}

function idsFrom(body) {
  let j;
  try {
    j = JSON.parse(body);
  } catch {
    return { ids: [], code: null };
  }
  const root = j.data || j;
  let models = root.models;
  if (models && !Array.isArray(models) && models.models) models = models.models;
  if (!Array.isArray(models)) return { ids: [], code: j.code, msg: j.msg };
  return { ids: models.map((m) => m.id || m.modelId).filter(Boolean), code: j.code, msg: j.msg };
}

(async () => {
  const bases = [
    ["copilot.tencent.com", "https://copilot.tencent.com/v3/config"],
    ["www.codebuddy.ai", "https://www.codebuddy.ai/v3/config"],
    ["www.codebuddy.cn", "https://www.codebuddy.cn/v3/config"],
  ];
  for (const acc of store.accounts) {
    console.log("===", acc.label, acc.site, "===");
    for (const [name, url] of bases) {
      try {
        const { status, body } = await req(url, acc.bearerToken);
        const { ids, code, msg } = idsFrom(body);
        const hy = ids.filter((x) => /hy4/i.test(x));
        console.log(
          name,
          "HTTP",
          status,
          "apiCode",
          code,
          "count",
          ids.length,
          "hy4",
          hy.join(",") || "(none)"
        );
        if (ids.length > 0 && ids.length <= 5) console.log("  ids:", ids.join(", "));
      } catch (e) {
        console.log(name, "ERR", e.message);
      }
    }
  }
})();
