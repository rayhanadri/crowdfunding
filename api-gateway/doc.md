sudo docker build -t gcr.io/crowdfunding-460613/api-gateway .

sudo docker push gcr.io/crowdfunding-460613/api-gateway

gcloud container images list --repository=gcr.io/crowdfunding-460613

gcloud run deploy api-gateway \
  --image gcr.io/crowdfunding-460613/api-gateway \
  --platform managed \
  --region asia-southeast2 \
  --allow-unauthenticated \
  --port 8080

 https://api-gateway-273575294549.asia-southeast2.run.app