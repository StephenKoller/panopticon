<template>
    <div class="task-list">
        <h1>Tasks</h1>
        <div class="task-controls">
            <input 
                v-model="newTaskTitle" 
                @keyup.enter="createTask"
                placeholder="Add new task..."
                type="text"
            >
            <button @click="createTask">Add Task</button>
        </div>
        <ul class="tasks">
            <li v-for="task in tasks" :key="task.id" class="task-item">
                <input 
                    type="checkbox" 
                    :checked="task.completed"
                    @change="toggleTask(task)"
                >
                <span :class="{ completed: task.completed }">{{ task.title }}</span>
                <button @click="deleteTask(task.id)" class="delete-btn">Delete</button>
            </li>
        </ul>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'

interface Task {
    id: number
    title: string
    completed: boolean
}

const tasks = ref<Task[]>([])
const newTaskTitle = ref('')

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

async function fetchTasks() {
    try {
        const response = await axios.get(`${apiUrl}/tasks`)
        tasks.value = response.data
    } catch (error) {
        console.error('Error fetching tasks:', error)
    }
}

async function createTask() {
    if (!newTaskTitle.value.trim()) return

    try {
        const response = await axios.post(`${apiUrl}/tasks`, {
            title: newTaskTitle.value,
            completed: false
        })
        tasks.value.push(response.data)
        newTaskTitle.value = ''
    } catch (error) {
        console.error('Error creating task:', error)
    }
}

async function toggleTask(task: Task) {
    try {
        await axios.put(`${apiUrl}/tasks/${task.id}`, {
            ...task,
            completed: !task.completed
        })
        task.completed = !task.completed
    } catch (error) {
        console.error('Error updating task:', error)
    }
}

async function deleteTask(id: number) {
    try {
        await axios.delete(`${apiUrl}/tasks/${id}`)
        tasks.value = tasks.value.filter(task => task.id !== id)
    } catch (error) {
        console.error('Error deleting task:', error)
    }
}

onMounted(fetchTasks)
</script>

<style scoped>
.task-list {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
}

.task-controls {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
}

.task-controls input {
    flex: 1;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
}

button {
    padding: 8px 16px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

button:hover {
    background-color: #45a049;
}

.tasks {
    list-style: none;
    padding: 0;
}

.task-item {
    display: flex;
    align-items: center;
    padding: 10px;
    border-bottom: 1px solid #eee;
    gap: 10px;
}

.task-item span {
    flex: 1;
}

.completed {
    text-decoration: line-through;
    color: #888;
}

.delete-btn {
    background-color: #f44336;
}

.delete-btn:hover {
    background-color: #da190b;
}
</style>
