<script lang="ts">
    import { onMount } from 'svelte';
    import Fa from 'svelte-fa/src/fa.svelte';
    import { faArrowsLeftRightToLine, faCaretDown, faCaretUp } from '@fortawesome/free-solid-svg-icons';

    enum AlarmState { Unset, Set, Running, Finished };

    let alarmState = AlarmState.Unset;
    let alarmTime = new Date();
    $: alarmHour = alarmTime.getHours();
    $: alarmMinute = alarmTime.getMinutes();
    $: alarmSecond = alarmTime.getSeconds();

    let timeDiffText : string = "";

    onMount(() => {
		  const interval = setInterval(() => { 
        const currTime: Date = new Date();
        switch (alarmState) {
          case AlarmState.Unset:
            // Update the clock to show the current time if no alarm is set
            alarmTime = currTime;
            break;
          case AlarmState.Running:
            // Update the time until text until alarm
            const [ hours, minutes, seconds ] = getTimeDiff(currTime.getTime(), alarmTime.getTime());
            if (hours == 0 && minutes == 0 && seconds == 0) {
              alarmState = AlarmState.Finished;
              timeDiffText = "";
            } else {
              timeDiffText = `${formatTime(hours)}:${formatTime(minutes)}:${formatTime(seconds)}`;
            }
            break;
          case AlarmState.Finished:
            playFinishAnim();
            break;
        }
      }, 1000);
		  return () => { clearInterval(interval); };
	  });

    const getTimeDiff = (currTimeMs: number, futureTimeMs: number): number[] => {
      const timeDiffMs: number = futureTimeMs - currTimeMs;
      const hours: number = Math.floor(timeDiffMs / 3600000);
      const minutes: number = Math.floor((timeDiffMs % 3600000) / 60000);
      const seconds: number = Math.floor((timeDiffMs % 60000) / 1000);
      return [ hours, minutes, seconds ];
    };

    let colorIndex = 0;
    const playFinishAnim = () => {
      const colors: string[] = ["#848C8E", "#435058", "#BFB7B6", "#8E4A49", "#5F464B"];
      document.body.style.backgroundColor = colors[colorIndex];
      colorIndex = (colorIndex + 1) % colors.length;
    };

    const formatTime = (n: number): string => (n <= 9) ? "0" + n : n.toString();

    const increaseHour = () => alarmHour = (alarmHour < 23) ? alarmHour + 1 : 0;
    const increaseMinute = () => alarmMinute = (alarmMinute < 59) ? alarmMinute + 1 : 0;
    const decreaseHour = () => alarmHour = (alarmHour == 0) ? 23 : alarmHour - 1;
    const decreaseMinute = () => alarmMinute = (alarmMinute == 0) ? 59 : alarmMinute - 1;

    const setAlarm = () => {
      alarmState = AlarmState.Set;
      alarmSecond = 0;
    };

    const startAlarm = () => {
      const currTime: Date = new Date();
      alarmTime.setHours(alarmHour, alarmMinute, alarmSecond);
      if (alarmTime.getTime() < currTime.getTime()) {
        alert("Invalid alarm time: Time has already passed!");
      }
      alarmState = AlarmState.Running;
    };

    const stopAlarm = () => {
      alarmState = AlarmState.Unset;
      alarmTime = new Date();
      document.body.style.backgroundColor = "black";
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
        <h1>:</h1>
        <h1 class="alarm-seconds">{formatTime(alarmSecond)}</h1>
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
          {:else if alarmState == AlarmState.Running} {timeDiffText} left
          {:else if alarmState == AlarmState.Finished} Alarm Finished!
          {/if}
    </pre>
    <div class="alarm-buttons">
        {#if alarmState == AlarmState.Unset}
          <button on:click={setAlarm}>Set Alarm</button>
        {:else if alarmState == AlarmState.Set}
          <button on:click={startAlarm}>Start Alarm</button>
        {:else if alarmState == AlarmState.Running}
          <button on:click={stopAlarm}>Stop Alarm</button>
        {:else if alarmState == AlarmState.Finished}
          <button on:click={stopAlarm}>Reset Alarm</button>
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