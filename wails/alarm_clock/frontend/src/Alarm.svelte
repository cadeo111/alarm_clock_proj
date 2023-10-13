<script lang="ts">
    import { onMount } from 'svelte';
    import Fa from 'svelte-fa/src/fa.svelte';
    import { faCaretDown, faCaretUp } from '@fortawesome/free-solid-svg-icons';

    let time: Date = new Date();

    enum AlarmState { Unset, Set, Running, Stopped };

    $: alarmState = AlarmState.Unset;
    $: alarmHour = time.getHours();
    $: alarmMinute = time.getMinutes();
    $: alarmSeconds = time.getSeconds();

    onMount(() => {
		  const interval = setInterval(() => { 
        if (alarmState == AlarmState.Unset) {
          time = new Date();
        }
      }, 1000);
		  return () => { clearInterval(interval); };
	  });

    const formatTime = (n: number): string => (n <= 9) ? "0" + n : n.toString();

    const increaseHour = () => alarmHour = (alarmHour < 23) ? alarmHour + 1 : 0;
    const increaseMinute = () => alarmMinute = (alarmMinute < 59) ? alarmMinute + 1 : 0;
    const decreaseHour = () => alarmHour = (alarmHour == 0) ? 23 : alarmHour - 1;
    const decreaseMinute = () => alarmMinute = (alarmMinute == 0) ? 59 : alarmMinute - 1;

    const setAlarm = (e) => {
      alarmState = AlarmState.Set;
      console.log("Alarm Set");
    };

    const startAlarm = (e) => {
      alarmState = AlarmState.Running;
      console.log("Alarm Running");
    };

    const stopAlarm = (e) => {
      alarmState = AlarmState.Stopped;
      console.log("Alarm Stopped");
    };

    const resetAlarm = (e) => {
      alarmState = AlarmState.Unset;
      console.log("Alarm Unset");
      time = new Date();
    };
</script>

<div>
    {#if alarmState == AlarmState.Set}
    <div class="alarm-setter">
        <button on:click={increaseHour}>
          <Fa icon={faCaretUp} size="7x" />
        </button>
        <button on:click={increaseMinute}>
          <Fa icon={faCaretUp} size="7x" />
        </button>
    </div>
    {/if}
    <div class="alarm-clock">
        <h1 class="alarm-hours">{formatTime(alarmHour)}</h1>
        <h1>:</h1>
        <h1 class="alarm-minutes">{formatTime(alarmMinute)}</h1>
    </div>
    {#if alarmState == AlarmState.Set}
    <div class="alarm-setter">
        <button on:click={decreaseHour}>
          <Fa icon={faCaretDown} size="7x" />
        </button>
        <button on:click={decreaseMinute}>
          <Fa icon={faCaretDown} size="7x" />
        </button>
    </div>
    {/if}
    <pre> {#if alarmState == AlarmState.Unset} Current Time
          {:else if alarmState == AlarmState.Set} Press Start
          {:else} Alarm Running!
          {/if}
    </pre>
    <div class="alarm-buttons">
        {#if alarmState == AlarmState.Unset}
          <button on:click={setAlarm}>Set Alarm</button>
        {:else if alarmState == AlarmState.Set}
          <button on:click={startAlarm}>Start</button>
          <button on:click={resetAlarm}>Reset</button>
        {:else}
          <button on:click={stopAlarm}>Stop</button>
          <button on:click={resetAlarm}>Reset</button>
        {/if}
    </div>
</div>

<style>
  .alarm-setter, .alarm-clock {
    display: flex;
    justify-content: center;
    align-items: center;
    margin: auto;
  }

  .alarm-setter > button {
    padding-left: 1rem;
    padding-right: 1rem;
  }

  .alarm-clock > h1 {
    color: rgba(250, 249, 246, 0.2);
    font-size: 7rem;
    height: 50%;
  }

  .alarm-buttons {
    padding-top: 2rem;
    display: inline-block;
    margin: auto;
    width: 50%;
  }
</style>