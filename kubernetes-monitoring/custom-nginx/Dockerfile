FROM nginx:1.24

# Создаем метрики
RUN echo "<html><body><h1>metrics: 42</h1></body></html>" > /usr/share/nginx/html/metrics.html

# Добавляем location /metrics
RUN sed -i '/location \/ {/a \\n    location /metrics {\n        alias /usr/share/nginx/html/;\n        index metrics.html;\n    }' /etc/nginx/conf.d/default.conf
