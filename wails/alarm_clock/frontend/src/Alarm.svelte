<script lang="ts">
    import {onMount} from 'svelte';
    import Fa from 'svelte-fa/src/fa.svelte';
    import {faCaretDown, faCaretUp} from '@fortawesome/free-solid-svg-icons';
    import {PlayAudio} from '../wailsjs/go/main/App.js'
    import {SnackbarContainer} from "attractions";
    import {StopAudio} from "../wailsjs/go/main/App";


    enum AlarmState { Unset, Set, Running, Finished };

    let alarmState = AlarmState.Unset;
    let alarmTime = new Date();

    let currTime = new Date();

    $: alarmHour = alarmTime.getHours();
    $: alarmMinute = alarmTime.getMinutes();
    $: alarmSecond = alarmTime.getSeconds();

    $: currHour = currTime.getHours();
    $: currMinute = currTime.getMinutes();
    $: currSecond = currTime.getSeconds();

    let timeDiffText: string = "";
    let playing = false;

    function play() {
        if (playing == false) {
            playing = true;
            PlayAudio(0);
        }
    }

    function stopPlay() {
        if (playing == true) {
            playing = false;
            StopAudio();
        }
    }

    onMount(() => {
        const interval = setInterval(() => {
            currTime = new Date();
            switch (alarmState) {
                case AlarmState.Unset:
                    stopPlay();
                    // Update the clock to show the current time if no alarm is set
                    alarmTime = currTime;
                    break;
                case AlarmState.Running:
                    stopPlay();
                    // Update the time until text until alarm
                    const [hours, minutes, seconds] = getTimeDiff(currTime.getTime(), alarmTime.getTime());
                    if (hours == 0 && minutes == 0 && seconds == 0) {
                        alarmState = AlarmState.Finished;
                        timeDiffText = "";
                    } else {
                        timeDiffText = `${formatTime(hours)}:${formatTime(minutes)}:${formatTime(seconds)}`;
                    }
                    break;
                case AlarmState.Finished:
                    play();
                    playFinishAnim();

                    break;
            }
        }, 1000);
        return () => {
            clearInterval(interval);
        };
    });

    const getTimeDiff = (currTimeMs: number, futureTimeMs: number): number[] => {
        const timeDiffMs: number = futureTimeMs - currTimeMs;
        const hours: number = Math.floor(timeDiffMs / 3600000);
        const minutes: number = Math.floor((timeDiffMs % 3600000) / 60000);
        const seconds: number = Math.floor((timeDiffMs % 60000) / 1000);
        return [hours, minutes, seconds];
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

    const startAlarm = (alertFunc: (msg: string) => void) => {
        const currTime: Date = new Date();
        alarmTime.setHours(alarmHour, alarmMinute, alarmSecond);
        console.log(alarmTime)
        if (alarmTime.getTime() < currTime.getTime()) {
            alarmTime.setDate(alarmTime.getDate() + 1);
        }
        // if (alarmTime.getTime() < currTime.getTime()) {
        //     console.log("Invalid alarm time: Time has already passed!")
        //     alertFunc("Invalid alarm time: Time has already passed!");
        //     return;
        // }
        alarmState = AlarmState.Running;
    };

    const stopAlarm = () => {
        alarmState = AlarmState.Unset;
        alarmTime = new Date();
        document.body.style.backgroundColor = "black";
    };
</script>


<div>

    {#if alarmState === AlarmState.Set}
        <div class="alarm-setter">
            <button on:click={()=>{
                console.log("hello");
                increaseHour()
            }}>
                <Fa icon={faCaretUp} size="7x"/>
            </button>
            <button on:click={increaseMinute}>
                <Fa icon={faCaretUp} size="7x"/>
            </button>
        </div>
    {/if}
    {#if (alarmState === AlarmState.Set)}
        <div class="alarm-clock">
            <h1 class="alarm-hours">{formatTime(alarmHour)}</h1>
            <h1>:</h1>
            <h1 class="alarm-minutes">{formatTime(alarmMinute)}</h1>
        </div>
    {:else }
        <div class="alarm-clock">
            <h1 class="alarm-hours">{formatTime(currHour)}</h1>
            <h1>:</h1>
            <h1 class="alarm-minutes">{formatTime(currMinute)}</h1>
            <h1>:</h1>
            <h1 class="alarm-seconds">{formatTime(currSecond)}</h1>
        </div>
    {/if}
    {#if alarmState === AlarmState.Set}
        <div class="alarm-setter">
            <button on:click={decreaseHour}>
                <Fa icon={faCaretDown} size="7x"/>
            </button>
            <button on:click={decreaseMinute}>
                <Fa icon={faCaretDown} size="7x"/>
            </button>
        </div>
    {/if}


    {#if alarmState === AlarmState.Set}
        <div class="snackbar-holder">
            <SnackbarContainer let:showSnackbar>
                <div class="set-controls">
                    <button class="button"  on:click={()=>{
               startAlarm((msg)=>showSnackbar({ props: { text: msg, class:"snackbar" } }))
            }}>Set Alarm
                    </button>

                </div>
            </SnackbarContainer>
        </div>
    {:else if alarmState === AlarmState.Running}
        <div class="alarm-text-group">
        <span class="alarm-time">{formatTime(alarmHour)}:{formatTime(alarmMinute)}</span>
        <span class="alarm-left">{timeDiffText} left</span>
        </div>
    {:else if alarmState === AlarmState.Finished}
        <pre>Alarm Finished!</pre>
    {/if}
    <div class="alarm-buttons">
        {#if alarmState === AlarmState.Unset}
            <button class="button"  on:click={setAlarm}>Set Alarm</button>
        {:else if alarmState === AlarmState.Set}
            <!--            <button on:click={startAlarm}>Start Alarm</button>-->
        {:else if alarmState === AlarmState.Running}
            <button class="button"  on:click={stopAlarm}>Stop Alarm</button>
        {:else if alarmState === AlarmState.Finished}
            <button class="button"  on:click={stopAlarm}>Reset Alarm</button>
        {/if}
    </div>


</div>


<style>

    .alarm-text-group{
        padding: 20px;
    }
    .alarm-left{
        padding: 4px 8px;
        background-color: #bd402b;
        border-radius: 5px;
    }
    .alarm-time{
        padding: 4px 8px;
        background-color: darkslateblue;
        border-radius: 5px;
    }
    .snackbar-holder {
        z-index: 3;
        position: absolute;
        top: 0;
        left: 0;
        bottom: 0;
        right: 0;
        /*height: 100vh;*/
        /*width: 100vw;*/
    }

    .set-controls {
        z-index: 2;
        position: absolute;
        bottom: 64px;
        right: 64px;
    }

    .button {
        padding: 24px;
        border: none;
        border-radius: 15px;
        font-size: 18px;
    }

    .set-controls .button {
        background-color: forestgreen;
        color: white;
    }

    .alarm-setter, .alarm-clock {
        display: flex;
        justify-content: center;
        align-items: center;
        margin: auto;
    }

    .alarm-setter > button {
        z-index: 4;
        padding-left: 1rem;
        padding-right: 1rem;
        margin-left: 32px;
    }

    .alarm-clock > h1 {
        color: rgba(250, 249, 246, 0.2);
        font-size: 7rem;
        height: 50%;
    }

    /*.alarm-buttons {*/
    /*    !*padding-top: 2rem;*!*/
    /*    !*display: inline-block;*!*/
    /*    !*margin: auto;*!*/
    /*    !*width: 50%;*!*/
    /*}*/
</style>