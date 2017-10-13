# redis_unique_queue

使用redis lua script 操作list + set数据结构, 构建redis的去重队列, 这样既能保证FIFO，又能保证去重.
