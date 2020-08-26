interface Event {
    type: 'EVT_TEXT' | 'EVT_JSON';
    data: string;
}

interface EventCallback {
    (event: Event): void;
}

interface ErrorCallback {
    (err: Error): void;
}

export default class EventStream {
    readonly target: string;
    readonly userId: string;
    private closed: boolean;

    constructor(target: string, userId: string) {
        this.target = target;
        this.userId = userId;
        this.closed = true;
    }

    subscribe(onEvent: EventCallback) {
        this.closed = false;
        const onError = (err: Error) => {
            if (!this.closed) {
                console.log(`subscribe event failure, retrying: ${err}`);
                window.setTimeout(() => {
                    this.subscribeEvent(onEvent, onError);
                }, 5000);
            }
        };
        this.subscribeEvent(onEvent, onError);
    }

    subscribeEvent(onEvent: EventCallback, onError: ErrorCallback) {
        const connectUrl = `http://${this.target}/sims/hub/connect`;
        const heartbeatUrl = `http://${this.target}/sims/hub/heartbeat`;
        const body = JSON.stringify({
            header: {
                user_id: this.userId,
            },
        });
        let ws: WebSocket;
        let beat = 0;

        function heartbeat() {
            window
                .fetch(heartbeatUrl, {method: 'POST', body: body})
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`event stream heartbeat: ${response.status}`);
                    }
                })
                .catch(err => {
                    if (ws) {
                        ws.close();
                    }
                    onError(err);
                });
        }

        window
            .fetch(connectUrl, {method: 'POST', body: body})
            .then(response => {
                if (!response.ok) {
                    throw new Error(`event stream connect: ${response.status}`);
                }
                return response.json();
            })
            .then(() => {
                ws = new WebSocket(`ws://${this.target}/sims/hub/events`);
                ws.onopen = () => {
                    ws.send(
                        JSON.stringify({
                            header: {
                                user_id: this.userId,
                                request_id: Math.floor(Date.now() / 1000).toString(),
                            },
                        }),
                    );
                    beat = window.setInterval(heartbeat, 5000);
                };
                ws.onmessage = ev => {
                    if (ev.data === '{}') return;
                    const event: Event = JSON.parse(ev.data);
                    event.data = atob(event.data);
                    onEvent(event);
                };
                ws.onclose = ev => {
                    window.clearInterval(beat);
                    onError(
                        new Error(
                            `event stream websocket closed: ${ev.code} (clean=${ev.wasClean})`,
                        ),
                    );
                };
                ws.onerror = ev => {
                    window.clearInterval(beat);
                    onError(new Error(`event stream websocket error: ${ev}`));
                };
            })
            .catch(err => {
                if (ws) {
                    ws.close();
                }
                onError(err);
            });
    }

    close() {
        this.closed = true;
        const disconnectUrl = `http://${this.target}/sims/hub/disconnect`;
        const body = JSON.stringify({
            header: {
                user_id: this.userId,
            },
        });
        window
            .fetch(disconnectUrl, {method: 'POST', body: body})
            .then(response => {
                if (!response.ok) {
                    throw new Error(`event stream disconnect: ${response.status}`);
                }
            })
            .catch(err => {
                console.log(err);
            });
    }
}
