addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
});

async function handleRequest(request) {
    if (request.url === 'https://updates.batect.dev/v1/latest') {
        return await fetch(new Request('https://storage.cloud.google.com/batect-updates-prod-public/v1/latest.json'));
    } else if (request.url === 'https://updates.batect.dev/ping') {
        return new Response('pong', { status: 200 });
    } else {
        return new Response('Not found', { status: 404 });
    }
}
