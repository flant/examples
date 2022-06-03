A Node.js application used to demonstrate how you can build your app and deploy it
to Kubernetes using [werf CI/CD tool](https://werf.io/).

Each directory corresponds to the state of our Node.js application when we
implement specific features in it:

* [01_basic_app](01_basic_app/) is the initial, basic state of our app;
* [020_assets](020_assets/) — when our app serves static assets;
* [030_db](030_db/) — when we connect to MySQL database;
* [040_s3](040_s3/) — when we use S3 storage for app's files.

P.S. This code is heavily based on the relevant werf guide:
[Node.js → Real-world apps](https://werf.io/guides/nodejs/200_real_apps.html).
