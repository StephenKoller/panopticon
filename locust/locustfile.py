from locust import HttpUser, task, between
from faker import Faker
import random
import json

fake = Faker()

class TaskUser(HttpUser):
    wait_time = between(1, 5)  # Wait between 1-5 seconds between tasks
    tasks_created = []  # Keep track of created task IDs

    def on_start(self):
        """Initialize user session"""
        self.client.headers = {'Content-Type': 'application/json'}

    @task(3)
    def get_tasks(self):
        """Get all tasks - higher weight as this is most common operation"""
        self.client.get("/tasks")

    @task(2)
    def create_task(self):
        """Create a new task"""
        task_data = {
            "title": fake.sentence(nb_words=6)[:-1],  # Remove trailing period
            "completed": False
        }
        
        response = self.client.post("/tasks", json=task_data)
        if response.status_code == 201:
            task = response.json()
            self.tasks_created.append(task["id"])

    @task(1)
    def update_task(self):
        """Update a task's status"""
        if not self.tasks_created:
            return

        task_id = random.choice(self.tasks_created)
        task_data = {
            "title": fake.sentence(nb_words=6)[:-1],
            "completed": random.choice([True, False])
        }
        
        self.client.put(f"/tasks/{task_id}", json=task_data)

    @task(1)
    def delete_task(self):
        """Delete a task"""
        if not self.tasks_created:
            return

        task_id = random.choice(self.tasks_created)
        response = self.client.delete(f"/tasks/{task_id}")
        
        if response.status_code == 200:
            self.tasks_created.remove(task_id)
