package org.rockets.cli_app.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.rockets.cli_app.components.Attachment;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class AttachmentService {
    private String baseUrl = "http://localhost:8080";
    HttpClient client = HttpClient.newHttpClient();
    ObjectMapper objectMapper = new ObjectMapper();

    public List<Attachment> getAttachments() {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/attachments"))
                .GET()
                .build();
        List<Attachment> attachments = new ArrayList<>();

        try {
            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
            String responseBody = response.body();
            ApiResponse<List<Attachment>> apiResponse = objectMapper.readValue(responseBody, ApiResponse.class);
            return apiResponse.getData();
        } catch (Error | IOException | InterruptedException e) {
            System.err.println(e.getMessage());
        }

        return attachments;
    }

    public Attachment createAttachment(Attachment attachment) {
        try {
            String requestBody = objectMapper.writeValueAsString(attachment);
            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(baseUrl + "/attachments"))
                    .header("Content-Type", "application/json")
                    .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                    .build();

            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
            ApiResponse<Attachment> apiResponse = objectMapper.readValue(response.body(), ApiResponse.class);
            return apiResponse.getData();
        } catch (IOException | InterruptedException e) {
            System.err.println(e.getMessage());
        }

        return null;
    }

    public Attachment updateAttachmentById(String id, Attachment attachment) {
        try {
            String requestBody = objectMapper.writeValueAsString(attachment);
            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(baseUrl + "/attachments/" + id))
                    .header("Content-Type", "application/json")
                    .method("PATCH", HttpRequest.BodyPublishers.ofString(requestBody))
                    .build();

            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
            ApiResponse<Attachment> apiResponse = objectMapper.readValue(response.body(), ApiResponse.class);
            return apiResponse.getData();
        } catch (IOException | InterruptedException e) {
            System.err.println(e.getMessage());
        }

        return null;
    }

    public void deleteAttachmentById(String id) {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/attachments/" + id))
                .DELETE()
                .build();

        try {
            client.send(request, HttpResponse.BodyHandlers.discarding());
        } catch (IOException | InterruptedException e) {
            System.err.println(e.getMessage());
        }
    }
}
