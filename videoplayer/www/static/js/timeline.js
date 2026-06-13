class D3Timeline {
    constructor(timelineId) {
        this.timelineElement = d3.select(`#${timelineId}`);
        this.data = [];  // Данные о видео
        this.width = this.timelineElement.node().clientWidth;
        this.height = this.timelineElement.node().clientHeight;
    }

    testData = [
        {Cam:{
            Name:"cam1", 
            files:[
            {StartVideo:"12:12:43",LengthVideo:"20",Link:"./cam1/2025-01-20/12_12_43.mp4"},
            {StartVideo:"12:13:04",LengthVideo:"20",Link:"./cam1/2025-01-20/12_13_04.mp4"},
            {StartVideo:"12:13:24",LengthVideo:"20",Link:"./cam1/2025-01-20/12_13_24.mp4"},
            {StartVideo:"12:13:45",LengthVideo:"20",Link:"./cam1/2025-01-20/12_13_45.mp4"},
            {StartVideo:"12:14:26",LengthVideo:"20",Link:"./cam1/2025-01-20/12_14_26.mp4"},
            {StartVideo:"12:14:47",LengthVideo:"20",Link:"./cam1/2025-01-20/12_14_47.mp4"},
            {StartVideo:"12:15:28",LengthVideo:"20",Link:"./cam1/2025-01-20/12_15_28.mp4"},
            {StartVideo:"12:15:49",LengthVideo:"20",Link:"./cam1/2025-01-20/12_15_49.mp4"},
            {StartVideo:"12:09:40",LengthVideo:"20",Link:"./cam1/2025-01-20/12_09_40.mp4"},
            {StartVideo:"12:10:01",LengthVideo:"20",Link:"./cam1/2025-01-20/12_10_01.mp4"}
            ]}
        },
        {Cam:{
            Name:"cam2", 
            files:[
            {StartVideo:"12:13:43",LengthVideo:"20",Link:"./cam1/2025-01-20/12_12_43.mp4"},
            {StartVideo:"12:14:04",LengthVideo:"20",Link:"./cam1/2025-01-20/12_13_04.mp4"},
            {StartVideo:"12:14:24",LengthVideo:"20",Link:"./cam1/2025-01-20/12_13_24.mp4"},
            {StartVideo:"12:14:45",LengthVideo:"20",Link:"./cam1/2025-01-20/12_13_45.mp4"},
            {StartVideo:"12:15:26",LengthVideo:"20",Link:"./cam1/2025-01-20/12_14_26.mp4"},
            {StartVideo:"12:15:47",LengthVideo:"20",Link:"./cam1/2025-01-20/12_14_47.mp4"},
            {StartVideo:"12:16:28",LengthVideo:"20",Link:"./cam1/2025-01-20/12_15_28.mp4"},
            {StartVideo:"12:16:49",LengthVideo:"20",Link:"./cam1/2025-01-20/12_15_49.mp4"},
            {StartVideo:"12:09:40",LengthVideo:"20",Link:"./cam1/2025-01-20/12_09_40.mp4"},
            {StartVideo:"12:10:01",LengthVideo:"20",Link:"./cam1/2025-01-20/12_10_01.mp4"}
            ]}
        }
    ]

    // Функция для получения данных с сервера
    async GetData(date) {
        // Имитация GET-запроса (замените URL на реальный API-эндпоинт)
        //const response = await fetch(`https://your-api.com/getVideos?date=${date}`);
        //const videos = await response.json();

        videos=this.testData;

        this.data = videos;
        this.renderTimeline();  // Отображение данных на шкале времени
    }

    // Функция для установки и центрирования временной шкалы
    SetCurrentTime(currentTime) {
        const scale = d3.scaleLinear()
            .domain([0, currentTime * 2])  // Пример: временной интервал от 0 до текущего времени * 2
            .range([0, this.width]);

        this.timelineElement.selectAll('.current-time-indicator').remove();

        // Добавляем индикатор текущего времени (можно сделать, например, вертикальную линию)
        this.timelineElement.append('div')
            .attr('class', 'current-time-indicator')
            .style('position', 'absolute')
            .style('top', '0')
            .style('left', `${scale(currentTime)}px`)
            .style('width', '2px')
            .style('height', `${this.height}px`)
            .style('background-color', 'red');
    }

    // Функция для отображения данных на временной шкале
    renderTimeline() {
        const barHeight = this.height * 0.8;  // Высота прямоугольников

        // Создаем шкалу для ширины времени
        const timeScale = d3.scaleTime()
            .domain([new Date().setHours(0, 0, 0, 0), new Date().setHours(24, 0, 0, 0)])  // Шкала на одни сутки
            .range([0, this.width]);

        // Преобразуем время начала видео в объекты Date
        const parseTime = d3.timeParse("%H:%M:%S");
        this.data.forEach(d => {
            d.StartVideo = parseTime(d.StartVideo);
        });

        // Отображаем каждый элемент как прямоугольник
        const bars = this.timelineElement.selectAll('.video-bar')
            .data(this.data)
            .join('div')
            .attr('class', 'video-bar')
            .style('position', 'absolute')
            .style('top', `${(this.height - barHeight) / 2}px`)  // Центрируем по высоте
            .style('left', d => `${timeScale(d.StartVideo)}px`)  // Устанавливаем левую позицию в зависимости от StartVideo
            .style('width', d => `${timeScale(new Date(d.StartVideo.getTime() + d.LengthVideo * 1000)) - timeScale(d.StartVideo)}px`)  // Ширина прямоугольника
            .style('height', `${barHeight}px`);
    }
}