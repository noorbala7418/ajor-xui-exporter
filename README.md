# Ajor XUI Exporter

Expose metrics on `0.0.0.0:9153/metrics`.

## Metrics

- ajor_xui_total_users_count
- ajor_xui_user_enable
- ajor_xui_user_admin_state
- ajor_xui_user_total_traffic
- ajor_xui_user_remain_traffic
- ajor_xui_user_upload_traffic
- ajor_xui_user_download_traffic

## Environment Variables

- `XPANEL_URL`: X-UI panel address. like `http://localhost:54321`.
- `XPANEL_USERNAME`: X-UI Username like `admin`.
- `XPANEL_PASSWORD`: X-UI Password like `admin`.
- `APP_LOG_MODE`: Log level mode. options `info` or `debug`. Default is `info`
- `APP_SCRAPE_TIME`: Scrape Time in seconds. Default is 30 seconds.

### Run

#### Run on your system

```bash
docker run -d --name ajor-xui-exporter \
-p 9153:9153 \
--network host \
--restart always \
-e XPANEL_URL="http://localhost:54321" \
-e XPANEL_USERNAME=admin \
-e APP_SCRAPE_TIME=20 \
-e XPANEL_PASSWORD=admin ajor-xui-exporter:latest
```
