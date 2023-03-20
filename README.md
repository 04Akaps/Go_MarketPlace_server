# MarketPlace를 만드는 서버

특정 서명에 대한 복호화, 채팅 기능 및 추가적인 기능은 작성하면서 사용할 예정

DB는 두개 사용할 예정
- mysql, mongo -> dynamoDB로 변경도 가능

캐시는 Redis사용 예정

Jenkins로 배포하고, 이번에는 Docker가 아니라 EC2의 로드 밸런싱 활용 예정
- CloudFront도 가능하다면??

이런식이면 실제로 비용이 부과될 것이니.. 일단 최대한 로컬망으로 작업 후 이후 EC2연동
- 사양한 가능하다면 프리티어로;;


# 사용하는 패키지
Viper
Gomniauth
Paseto
...++


# env 관리

내부에 app.env로 설정하여 다음과 같은 설정을 실행

# 프로젝트 개요
MarketPlace를 하나 만들 예정
모든 Network에 대해서 블록을 계속 풀링하정건 불가능하니.. 해당 사이트에서 Launchpad를 하나 두어서
거기에서 발생하는 NFT들에 대해서 데이터를 추적하여 처리 할 예정

Launchpad는 ERC721A로 작성하고, Network는 Polygon 정도로 생각 중

추가적으로 Chatting기능도 넣을 예

