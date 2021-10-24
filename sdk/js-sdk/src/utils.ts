import { setTimeout } from 'timers/promises';

export async function sleep(ms: number) {
  await setTimeout(ms);
}
