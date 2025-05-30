// Конфигурация API
const API_BASE_URL = 'http://localhost:3002/api/forum';
const AUTH_API_URL = 'http://localhost:3001/api';

// Функции для работы с аутентификацией
async function register(email, username, password) {
    try {
        const response = await fetch(`${AUTH_API_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, username, password }),
        });
        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.message || 'Registration failed');
        }
        return data;
    } catch (error) {
        console.error('Registration error:', error);
        throw error;
    }
}

async function login(username, password) {
    try {
        console.log('Attempting login with username:', username);
        const response = await fetch(`${AUTH_API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });
        
        let data;
        const contentType = response.headers.get('content-type');
        
        if (!response.ok) {
            // Если ответ не успешный, пробуем получить текст ошибки
            const errorText = await response.text();
            console.error('Login failed with status:', response.status);
            console.error('Error response:', errorText);
            throw new Error(errorText || 'Login failed');
        }
        
        // Если ответ успешный, парсим JSON
        data = await response.json();
        console.log('Login successful, received data:', data);
        
        if (!data.user || !data.user.id) {
            console.error('Invalid server response: missing user data');
            throw new Error('Invalid server response: missing user data');
        }

        if (!data.token) {
            console.error('No token received from server');
            throw new Error('No token received from server');
        }

        console.log('Storing token and user data in localStorage');
        localStorage.setItem('token', data.token);
        localStorage.setItem('username', username);
        localStorage.setItem('user_id', data.user.id);
        localStorage.setItem('user_role', data.user.role);
        
        // Verify storage
        console.log('Verifying stored data:');
        console.log('Token:', localStorage.getItem('token'));
        console.log('Username:', localStorage.getItem('username'));
        console.log('User ID:', localStorage.getItem('user_id'));
        console.log('User Role:', localStorage.getItem('user_role'));
        
        updateAuthUI(true);
        return data;
    } catch (error) {
        console.error('Login error:', error);
        throw error;
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    localStorage.removeItem('user_id');
    localStorage.removeItem('user_role');
    updateAuthUI(false);
    showSection('home');
}

function updateAuthUI(isAuthenticated) {
    const authLinks = document.getElementById('auth-links');
    const userLinks = document.getElementById('user-links');
    const usernameDisplay = document.getElementById('username-display');
    
    if (isAuthenticated) {
        authLinks.classList.add('hidden');
        userLinks.classList.remove('hidden');
        usernameDisplay.textContent = localStorage.getItem('username');
    } else {
        authLinks.classList.remove('hidden');
        userLinks.classList.add('hidden');
        usernameDisplay.textContent = '';
    }
}

// Функции для работы с категориями
async function getCategories() {
    try {
        const response = await fetch(`${API_BASE_URL}/categories`);
        if (!response.ok) {
            throw new Error('Failed to fetch categories');
        }
        return await response.json();
    } catch (error) {
        console.error('Error fetching categories:', error);
        throw error;
    }
}

async function createCategory(name, description) {
    try {
        const response = await fetch(`${API_BASE_URL}/categories`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify({ name, description })
        });
        if (!response.ok) {
            throw new Error('Failed to create category');
        }
        return await response.json();
    } catch (error) {
        console.error('Error creating category:', error);
        throw error;
    }
}

async function deleteCategory(categoryId) {
    try {
        const response = await fetch(`${API_BASE_URL}/delete_category?id=${categoryId}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        if (!response.ok) {
            throw new Error('Failed to delete category');
        }
    } catch (error) {
        console.error('Error deleting category:', error);
        throw error;
    }
}

function updateCategoriesList(categories) {
    const categoriesList = document.getElementById('categoriesList');
    const postCategory = document.getElementById('postCategory');
    
    if (!categoriesList || !postCategory) return;
    
    categoriesList.innerHTML = '';
    postCategory.innerHTML = '<option value="">Select Category</option>';
    
    if (!categories || !Array.isArray(categories)) {
        categories = [];
    }
    
    const userRole = localStorage.getItem('user_role');
    const isAdmin = userRole === 'admin';
    
    categories.forEach(category => {
        // Добавляем категорию в список
        const categoryElement = document.createElement('div');
        categoryElement.className = 'category';
        categoryElement.innerHTML = `
            <div class="category-header">
                <h3 class="category-title">${category.name}</h3>
                ${isAdmin ? `
                <button class="delete-category-btn" onclick="deleteCategory(${category.id})" title="Delete category">
                    <span class="delete-icon">×</span>
                </button>
                ` : ''}
            </div>
            <p>${category.description}</p>
            <button onclick="loadPosts(${category.id})">View Posts</button>
        `;
        categoriesList.appendChild(categoryElement);
        
        // Добавляем категорию в селект для создания поста
        const option = document.createElement('option');
        option.value = category.id;
        option.textContent = category.name;
        postCategory.appendChild(option);
    });
}

// Функции для работы с постами
async function getPosts() {
    try {
        const response = await fetch(`${API_BASE_URL}/posts`);
        if (!response.ok) {
            throw new Error('Failed to fetch posts');
        }
        return await response.json();
    } catch (error) {
        console.error('Error fetching posts:', error);
        throw error;
    }
}

async function createPost(title, content, categoryId) {
    try {
        const response = await fetch(`${API_BASE_URL}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify({ title, content, categoryId })
        });
        if (!response.ok) {
            throw new Error('Failed to create post');
        }
        return await response.json();
    } catch (error) {
        console.error('Error creating post:', error);
        throw error;
    }
}

// Функции для работы с комментариями
async function getComments(postId) {
    try {
        const response = await fetch(`${API_BASE_URL}/comments?post_id=${postId}`);
        if (!response.ok) {
            throw new Error('Failed to fetch comments');
        }
        return await response.json();
    } catch (error) {
        console.error('Error fetching comments:', error);
        throw error;
    }
}

async function createComment(content, postId) {
    try {
        const response = await fetch(`${API_BASE_URL}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify({ content, postId })
        });
        if (!response.ok) {
            throw new Error('Failed to create comment');
        }
        return await response.json();
    } catch (error) {
        console.error('Error creating comment:', error);
        throw error;
    }
}

// Функции для удаления
async function deletePost(postId) {
    try {
        const response = await fetch(`${API_BASE_URL}/delete_post?id=${postId}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        if (!response.ok) {
            throw new Error('Failed to delete post');
        }
    } catch (error) {
        console.error('Error deleting post:', error);
        throw error;
    }
}

async function deleteComment(commentId) {
    try {
        const response = await fetch(`${API_BASE_URL}/delete_comment?id=${commentId}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
        });
        if (!response.ok) {
            throw new Error('Failed to delete comment');
        }
    } catch (error) {
        console.error('Error deleting comment:', error);
        throw error;
    }
}

// Функции для работы с чатом
let ws = null;

function connectToChat() {
    ws = new WebSocket('ws://localhost:3003/ws');
    
    ws.onopen = () => {
        console.log('Connected to chat');
    };
    
    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        displayMessage(message);
    };
    
    ws.onclose = () => {
        console.log('Disconnected from chat');
        setTimeout(connectToChat, 1000);
    };
}

function sendMessage(content) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
            content,
            user_id: parseInt(localStorage.getItem('user_id')),
            username: localStorage.getItem('username')
        }));
    }
}

function displayMessage(message) {
    const chatMessages = document.querySelector('.chat-messages');
    if (chatMessages) {
        const messageElement = document.createElement('div');
        messageElement.className = `message ${message.username === localStorage.getItem('username') ? 'own-message' : ''}`;
        messageElement.innerHTML = `
            <span class="username">${message.username}:</span>
            <p>${message.content}</p>
        `;
        chatMessages.appendChild(messageElement);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
}

// Функции для управления UI
function showElement(id) {
    document.getElementById(id).classList.remove('hidden');
}

function hideElement(id) {
    document.getElementById(id).classList.add('hidden');
}

function showAuthForms() {
    showElement('auth-forms');
    hideElement('categories');
    hideElement('posts');
    hideElement('chat');
}

function showMainContent() {
    hideElement('auth-forms');
    showElement('categories');
    showElement('posts');
    showElement('chat');
}

function showSection(sectionId) {
    // Скрываем все секции
    const sections = document.querySelectorAll('.content-section');
    sections.forEach(section => {
        section.classList.add('hidden');
    });

    // Показываем нужную секцию
    const targetSection = document.getElementById(sectionId);
    if (targetSection) {
        targetSection.classList.remove('hidden');
    }

    // Обновляем активную ссылку
    const navLinks = document.querySelectorAll('.nav-link');
    navLinks.forEach(link => {
        link.classList.remove('active');
        if (link.dataset.section === sectionId) {
            link.classList.add('active');
        }
    });

    // Загружаем данные в зависимости от раздела
    if (sectionId === 'categories') {
        getCategories().then(categories => {
            updateCategoriesList(categories);
        }).catch(error => {
            console.error('Error loading categories:', error);
        });
    } else if (sectionId === 'posts') {
        // Получаем текущую категорию из localStorage или используем первую доступную
        const currentCategoryId = localStorage.getItem('currentCategoryId');
        if (currentCategoryId) {
            getPosts(currentCategoryId).then(posts => {
                updatePostsList(posts);
            }).catch(error => {
                console.error('Error loading posts:', error);
            });
        } else {
            // Если нет выбранной категории, загружаем все категории и выбираем первую
            getCategories().then(categories => {
                if (categories && categories.length > 0) {
                    localStorage.setItem('currentCategoryId', categories[0].id);
                    getPosts(categories[0].id).then(posts => {
                        updatePostsList(posts);
                    }).catch(error => {
                        console.error('Error loading posts:', error);
                    });
                }
            }).catch(error => {
                console.error('Error loading categories:', error);
            });
        }
    }
}

// Навигация
document.addEventListener('DOMContentLoaded', () => {
    // Обработка навигации
    const navLinks = document.querySelectorAll('.nav-link');

    // Обработчики клика по навигации
    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const sectionId = link.dataset.section;
            showSection(sectionId);
        });
    });

    // Показываем домашнюю страницу по умолчанию
    showSection('home');

    // Проверяем авторизацию
    const token = localStorage.getItem('token');
    updateAuthUI(!!token);

    // Обработчик кнопки выхода
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', logout);
    }

    // Подключаемся к чату если авторизованы
    if (token) {
        connectToChat();
    }
    
    // Обработчики форм
    const registerForm = document.getElementById('registerForm');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            try {
                const email = document.getElementById('registerEmail').value;
                const username = document.getElementById('registerUsername').value;
                const password = document.getElementById('registerPassword').value;
                await register(email, username, password);
                alert('Registration successful! Please login.');
                showSection('login');
            } catch (error) {
                alert(error.message);
            }
        });
    }
    
    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            try {
                const username = document.getElementById('loginUsername').value;
                const password = document.getElementById('loginPassword').value;
                await login(username, password);
                showSection('categories');
                // Загружаем категории
                const categories = await getCategories();
                updateCategoriesList(categories);
            } catch (error) {
                alert(error.message);
            }
        });
    }
    
    const categoryForm = document.getElementById('categoryForm');
    if (categoryForm) {
        categoryForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            try {
                const name = document.getElementById('categoryName').value;
                const description = document.getElementById('categoryDescription').value;
                await createCategory(name, description);
                const categories = await getCategories();
                updateCategoriesList(categories);
                categoryForm.reset();
            } catch (error) {
                alert(error.message);
            }
        });
    }

    const postForm = document.getElementById('postForm');
    if (postForm) {
        postForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            try {
                const title = document.getElementById('postTitle').value;
                const content = document.getElementById('postContent').value;
                const categoryId = document.getElementById('postCategory').value;
                
                if (!categoryId) {
                    alert('Please select a category');
                    return;
                }
                
                await createPost(title, content, categoryId);
                postForm.reset();
                
                // Обновляем список постов после создания
                const posts = await getPosts();
                updatePostsList(posts);
                
                // Показываем раздел с постами
                showSection('posts');
            } catch (error) {
                alert(error.message);
            }
        });
    }
    
    // Обработчик формы чата
    const chatForm = document.querySelector('.chat-input');
    if (chatForm) {
        chatForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const input = chatForm.querySelector('input');
            const message = input.value.trim();
            if (message) {
                sendMessage(message);
                input.value = '';
            }
        });
    }
});

// Функции для обновления UI
async function loadPosts(categoryId) {
    try {
        localStorage.setItem('currentCategoryId', categoryId);
        const posts = await getPosts();
        updatePostsList(posts);
        showSection('posts');
    } catch (error) {
        console.error('Error loading posts:', error);
        // Если категория не найдена, показываем пустой список
        updatePostsList([]);
    }
}

function updatePostsList(posts) {
    const postsList = document.getElementById('postsList');
    if (!postsList) return;
    
    postsList.innerHTML = '';
    
    if (!posts || !Array.isArray(posts)) {
        posts = [];
    }
    
    if (posts.length === 0) {
        postsList.innerHTML = '<p class="no-posts">No posts in this category yet.</p>';
        return;
    }
    
    // Получаем текущую категорию
    const currentCategoryId = localStorage.getItem('currentCategoryId');
    let currentCategoryName = 'Unknown Category';
    
    // Получаем название категории из селекта
    const postCategory = document.getElementById('postCategory');
    if (postCategory) {
        const selectedOption = postCategory.querySelector(`option[value="${currentCategoryId}"]`);
        if (selectedOption) {
            currentCategoryName = selectedOption.textContent;
        }
    }
    
    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.className = 'post';
        postElement.innerHTML = `
            <div class="post-header">
                <h3 class="post-title">${post.title}</h3>
                <span class="post-category">Category: ${currentCategoryName}</span>
            </div>
            <p class="post-content">${post.content}</p>
            <div class="post-meta">
                <small>Posted by ${post.authorId || 'Anonymous'} on ${new Date(post.createdAt).toLocaleString()}</small>
            </div>
        `;
        postsList.appendChild(postElement);
    });
} 