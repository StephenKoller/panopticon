<template>
    <div id="app">
        <nav>
            <h1>Panopticon Demo</h1>
        </nav>
        <router-view></router-view>
    </div>
</template>

<style>
body {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
        Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    background-color: #f5f5f5;
}

#app {
    min-height: 100vh;
}

nav {
    background-color: #2c3e50;
    color: white;
    padding: 1rem;
    text-align: center;
}

nav h1 {
    margin: 0;
}
</style>
