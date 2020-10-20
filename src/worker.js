addEventListener('fetch', event => {
    event.respondWith(handleRequest(event.request))
});

async function handleRequest(request) {
    const originResponse = await getResponse(request);
    const modifiedResponse = new Response(originResponse.body, originResponse);

    modifiedResponse.headers.set('Content-Security-Policy', "default-src 'none'; frame-ancestors 'none'");
    modifiedResponse.headers.set('X-Frame-Options', 'DENY');
    modifiedResponse.headers.set('X-Content-Type-Options', 'nosniff');
    modifiedResponse.headers.set('Referrer-Policy', 'no-referrer');

    return modifiedResponse;
}

async function getResponse(request) {
    if (request.url === 'https://updates.batect.dev/v1/latest') {
        return await fetch(new Request('https://storage.googleapis.com/batect-updates-prod-public/v1/latest.json'));
    } else if (request.url === 'https://updates.batect.dev/ping') {
        return new Response('pong', { status: 200 });
    } else if (request.url === 'https://updates.batect.dev/') {
        return new Response('', { status: 200 });
    } else {
        return new Response('Not found', { status: 404 });
    }
}
